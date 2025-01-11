package app

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
	"github.com/patyhank/ptd/app/config"
	"github.com/patyhank/ptd/app/ent"
	"github.com/patyhank/ptd/app/ent/author"
	"github.com/patyhank/ptd/app/ent/message"
	"github.com/patyhank/ptd/app/ent/postinfo"
	"github.com/patyhank/ptd/core"
	"github.com/patyhank/ptd/core/event"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
	"sync"
)

type Instance struct {
	config config.Config

	discord      bot.Client
	currentPost  *ent.PostInfo // 目前瀏覽的貼文
	previousPost *ent.PostInfo // 上一個瀏覽的貼文，瀏覽下一篇時，於Discord新增已完成標籤
	db           *ent.Client
	*core.Client

	currentSearchIndex       int
	currentSearch            config.SearchConfig
	currentVariantIndex      int
	currentTitleVariantIndex int

	fetchedPostAID string
	fetchedPostURL string

	reFetchSignal chan bool

	authorCache sync.Map

	viewState ViewState
	sync.Mutex
}

func NewInstance(cfg config.Config) *Instance {
	client, err := disgo.New(cfg.Discord.Token,
		bot.WithCacheConfigOpts(cache.WithCaches(cache.FlagsAll)),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentsNonPrivileged)))
	if err != nil {
		log.Fatalf("error creating discord client: %v", err)
		return nil
	}
	var dbClient *ent.Client
	dbClient, err = ent.Open("sqlite3", "file:data.db?&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
		return nil
	}
	// Run the auto migration tool.
	if err := dbClient.Schema.Create(context.Background(), schema.WithDropIndex(true), schema.WithDropColumn(true)); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
		return nil
	}

	conn := core.NewConn(cfg.PTT.Connection.Host, cfg.PTT.Connection.HostOrigin)

	return &Instance{
		discord:       client,
		config:        cfg,
		db:            dbClient,
		Client:        conn,
		reFetchSignal: make(chan bool, 1),
	}
}

func (i *Instance) Close(ctx context.Context) error {
	if err := i.db.Close(); err != nil {
		return err
	}
	i.discord.Close(ctx)

	return nil
}

func (i *Instance) RegisterAccountHandler() {
	i.AddEventListeners(
		core.NewListenerFunc(func(c *core.EventClient, e *event.PressAnyKeyEvent) {
			i.SendReturn()
		}),
		core.NewListenerFunc(func(c *core.EventClient, e *event.BadLoginNotifyEvent) {
			c.SendMessage("n", true)
		}),
		core.NewListenerFunc(func(c *core.EventClient, e *event.DuplicateConnectionEvent) {
			c.SendMessage("n", true)
		}),
		core.NewListenerFunc(i.mainScreen),
		core.NewListenerFunc(i.viewPost),
		core.NewListenerFunc(i.comment),
		core.NewListenerFunc(i.postInfo),
	)
}
func (i *Instance) recordComment(c event.CommentData) {
	var emoji discord.Emoji

	switch c.Type {
	case event.CommentTypeUpVote:
		if i.currentSearch.Emoji.UpVote.Name != "" {
			emoji = i.currentSearch.Emoji.UpVote
		} else {
			emoji.Name = "👍"
		}
	case event.CommentTypeDownVote:
		if i.currentSearch.Emoji.DownVote.Name != "" {
			emoji = i.currentSearch.Emoji.DownVote
		} else {
			emoji.Name = "👎"
		}
	case event.CommentTypeReply:
		if i.currentSearch.Emoji.Reply.Name != "" {
			emoji = i.currentSearch.Emoji.Reply
		} else {
			emoji.Name = "↩️"
		}
	}

	data := fmt.Sprintf("%v **%v** : %v", emojiName(emoji), c.Author, c.Content)
	plainText := fmt.Sprintf("%s %s\t:\t%s\t%s", c.Type.String(), c.Author, c.Content, c.Time.Format("15:04"))
	hashByte := sha512.Sum512([]byte(plainText))
	hash := hex.EncodeToString(hashByte[:])

	exist, _ := i.currentPost.QueryMessages().Where(message.HashEQ(hash)).Exist(context.Background())
	if exist {
		return
	}
	var authorData *ent.Author

	authorAny, ok := i.authorCache.Load(c.Author)
	if !ok {
		var err error
		authorData, err = i.db.Author.Query().Where(author.AuthorIDEqualFold(c.Author)).First(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				authorData, err = i.db.Author.Create().SetAuthorID(c.Author).SetLastSeen(time.Now()).Save(context.Background())
				if err != nil {
					log.Warn(err)
				}
			} else {
				log.Warn(err)
			}
		} else {
			authorData, err = authorData.Update().SetLastSeen(time.Now()).Save(context.Background())
			if err != nil {
				log.Warn(err)
			}
		}
		i.authorCache.Store(c.Author, authorData)
	} else {
		authorData = authorAny.(*ent.Author)
	}

	err := i.db.Message.Create().SetContent(data).SetRawContent(plainText).SetAuthorID(authorData.ID).SetCreatedAt(time.Now()).SetHash(hash).SetParentPost(i.currentPost).Exec(context.Background())
	if err != nil {
		log.Warn(err)
	}
}

func (i *Instance) comment(c *core.EventClient, e *event.CommentDataEvent) {
	if i.currentPost == nil {
		return
	}
	for _, comment := range e.Comments {
		i.recordComment(comment)
	}
	return
}
func (i *Instance) postInfo(c *core.EventClient, e *event.PostInfoScreenEvent) {
	log.Infof("已擷取到文章資訊")
	i.fetchedPostAID = e.PostAID
	i.fetchedPostURL = e.PostURL
}
func (i *Instance) mainScreen(c *core.EventClient, e *event.MainScreenEvent) {
	log.Infof("已進入主畫面")
	i.currentSearchIndex++
	if i.currentSearchIndex >= len(i.config.Searches) {
		i.currentSearchIndex = 0
	}
	i.currentSearch = i.config.Searches[i.currentSearchIndex]

	i.Lock()
	if i.viewState > ViewStateBacking {
		log.Infof("Current view state is %v, skipping searching, ", i.viewState)
		i.Unlock()
		return
	}

	if i.viewState <= ViewStateBacking {
		log.Infof("Current view state is %v, changing to searching", i.viewState)
		i.viewState = ViewStateSearching
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
	defer cancelFunc()
	defer func() {
		log.Infof("Change state to ready viewing")
		i.viewState = ViewStateReadyViewing
		i.Unlock()
		go func() {
			i.SendRefresh()
			i.SendRefresh()
		}()
	}()

	c.PrepareWait()
	i.SendMultipleMessage([]string{"s", i.currentSearch.Board})
	c.Wait(ctx)
	i.SendReturn()
	time.Sleep(time.Second)
	customView, _ := i.db.PostInfo.Query().Where(postinfo.ForceViewExpireGTE(time.Now())).First(ctx)
	if customView != nil {
		if customView.Aid != "" {
			i.SendMessage("#")
			i.SendMessage(customView.Aid, true)
			i.viewState = ViewStateReadyViewing

			c.PrepareWait()
			i.SendRefresh()
			log.Infof("Change state to ready viewing, Current viewing custom AID: %s", customView.Aid)
			return
		}

		if i.currentTitleVariantIndex >= len(i.currentSearch.TitleSearchVariant) {
			i.currentTitleVariantIndex = 0
		}

		pattern := i.currentSearch.TitleSearchVariant[i.currentTitleVariantIndex]
		i.SendMultipleMessage(pattern.Keys())

		for _, keyword := range customView.SearchKeywords {
			i.SendMessage("/")
			i.SendMessage(keyword, true)
		}
		i.viewState = ViewStateReadyViewing

		i.currentTitleVariantIndex++

		log.Infof("Change state to ready viewing, Current viewing custom AID: %s", customView.Aid)
		return
	}

	if i.currentVariantIndex >= len(i.currentSearch.SearchVariant) {
		i.currentVariantIndex = 0
	}

	pattern := i.currentSearch.SearchVariant[i.currentVariantIndex]
	i.SendMultipleMessage(pattern.Keys())
	i.currentVariantIndex++

	i.SendMessage("\x1b[6~")
	i.SendMessage("\x1b[6~")
	i.viewState = ViewStateReadyViewing
	i.SendRefresh()
}

func (i *Instance) viewPost(c *core.EventClient, e *event.ListPostEvent) {
	if i.viewState != ViewStateReadyViewing {
		if i.viewState != ViewStateViewing {
			log.Warn("Current view state is not ready viewing, skipping ", i.viewState.String())
			time.Sleep(time.Second)
			i.SendRefresh()
		}
		return
	}
	if len(e.Posts) == 0 {
		log.Warn("No posts found, backing")
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	c.PrepareWait()
	i.SendRefresh()
	c.Wait(ctx)

	regex := regexp.MustCompile(i.currentSearch.PostMatchRegex)
	var posts []event.PostInfo

	for _, post := range e.Posts {
		if !regex.MatchString(post.Title) {
			continue
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		log.Warn("No posts found, backing")
		return
	}

	customView, _ := i.db.PostInfo.Query().Where(postinfo.ForceViewExpireGTE(time.Now())).First(ctx)
	lastPost := lo.LastOrEmpty(posts)
	var forceView bool
	if customView != nil {
		forceView = true
		if customView.Aid == "" && customView.SearchKeywords != nil {
			for _, post := range posts {
				match := regex.FindStringSubmatch(post.Title)
				if match != nil {
					postTitle := passingStringsToFmt(i.currentSearch.PostTitle, match)
					if postTitle == customView.Title {
						lastPost = post
						forceView = true

						_, err := customView.Update().ClearForceViewExpire().Save(ctx)

						if err != nil {
							log.Warn(err)
						}

						break
					}
				}
			}
		}
		if customView.Aid != "" {
			for _, post := range posts {
				if post.Cursor {
					lastPost = post

					err := customView.Update().ClearForceViewExpire().Exec(ctx)

					if err != nil {
						log.Warn(err)
					}
					break
				}
			}
		}
	}

	if lastPost.Title == "" {
		log.Infof("No posts found, backing")
		i.currentPost = nil
		i.Lock()
		i.viewState = ViewStateBacking
		i.Unlock()
		err := i.Backing(ctx)
		if err != nil {
			log.Errorf("error backing: %v", err)
		}
	}
	i.Lock()
	i.viewState = ViewStateViewing
	i.Unlock()

	postTitle := passingStringsToFmt(i.currentSearch.PostTitle, regex.FindStringSubmatch(lastPost.Title))
	c.PrepareWait()
	i.SendMessage("Q")
	i.SendMessage("\f")
	c.Wait(ctx)
	i.SendMessage("\f")
	c.PrepareWait()
	c.Wait(ctx)
	c.PrepareWait()
	i.SendReturn()
	c.Wait(ctx)
	c.PrepareWait()
	i.SendRefresh()
	c.Wait(ctx)

	contentStart := "───────────────────────────────────────"

	log.Infof("已擷取文章資訊, 文章: %s [%s]", lastPost.Title, postTitle)
	log.Infof("文章AID: %v, 文章網址: %v", i.fetchedPostAID, i.fetchedPostURL)

	pageContent := i.String()
	index := strings.Index(pageContent, contentStart)
	if index == -1 {
		i.Lock()
		i.viewState = ViewStateBacking
		i.Unlock()
		if err := i.Backing(ctx); err != nil {
			log.Errorf("error backing: %v", err)
		}
		return
	}
	senderRegex := regexp.MustCompile("(?:--\\s+\\n(?:[^\\n]*?\\n){0,6})?--\\s+\\n※ 發信站: ")

	pageContent = pageContent[index+len(contentStart):]
	{
		postEndIndex := senderRegex.FindStringIndex(pageContent)
		if postEndIndex != nil {
			pageContent = pageContent[:postEndIndex[0]-1]
		}
		lastIndex := strings.LastIndex(pageContent, "\n  瀏覽 第 ")
		if lastIndex != -1 {
			pageContent = pageContent[:lastIndex]
		}
	}
	for !senderRegex.MatchString(i.String()) {
		i.SendMessage("\x20")
		c.PrepareWait()
		i.SendMessage("\f")
		c.Wait(ctx)
		screen := i.String()
		postEndIndex := senderRegex.FindStringIndex(screen)
		if postEndIndex != nil {
			screen = screen[:postEndIndex[0]-1]
			break
		}

		lastIndex := strings.LastIndex(screen, "\n  瀏覽 第 ")
		if lastIndex != -1 {
			pageContent += screen[:lastIndex-1]
		}
	}

	contentChunk := SplitContentByLength(pageContent, 1800)

	c.PrepareWait()

	postInfo, err := i.db.PostInfo.Query().Where(postinfo.TitleEqualFold(postTitle)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		log.Errorf("error querying post info: %v", err)
	}

	if postInfo == nil {
		i.currentPost, err = i.db.PostInfo.Create().SetTitle(postTitle).SetURL(i.fetchedPostURL).SetAid(i.fetchedPostAID).Save(ctx)
		if err != nil {
			log.Errorf("error saving post info: %v", err)
		}

		i.createThreadAndPinnedMessage(contentChunk)
	} else {
		i.currentPost = postInfo

		if len(postInfo.ContentMessages) == 0 {
			i.createThreadAndPinnedMessage(contentChunk)
		}
	}

	if forceView {
		if i.currentSearch.ForumChannel != 0 {
			err := i.discord.Rest().RemoveOwnReaction(i.currentPost.ChannelID, i.currentPost.ChannelID, "✅")
			if err != nil {
				log.Debug(err)
			}
		}
	}
	if i.currentPost.PostContent != pageContent && i.currentPost.LastUpdated.Add(time.Minute*10).Before(time.Now()) {
		go func() {
			log.Infof("正在更新訊息內文，請忽略警告訊息")
			for j, str := range contentChunk {
				if j >= i.currentSearch.PreFillSize {
					break
				}
				message, err := i.discord.Rest().GetMessage(i.currentPost.ChannelID, i.currentPost.ContentMessages[j])
				if err != nil {
					log.Warn(err)
				}
				if strings.ReplaceAll(str, " ", "") != strings.ReplaceAll(message.Content, " ", "") {
					_, err = i.discord.Rest().UpdateMessage(i.currentPost.ChannelID, i.currentPost.ContentMessages[j], discord.NewMessageUpdateBuilder().SetContent(str).Build())
					if err != nil {
						log.Debug(err)
					}
				}
			}
			log.Infof("內文更新完成")
		}()

		i.currentPost, err = i.currentPost.Update().SetPostContent(pageContent).SetLastUpdated(time.Now()).Save(context.Background())
		if err != nil {
			log.Warn(err)
		}
	}

	i.SendMessage("\u001B[B")
	i.SendMessage("\u001B[A")
	c.PrepareWait()
	i.SendRefresh()
	c.Wait(ctx)

	for !strings.Contains(i.String(), "頁 (100%)") {
		i.SendMessage("\x20")
		c.PrepareWait()
		i.SendRefresh()
		c.Wait(ctx)
		i.CheckAndSendMessages()
		if strings.Contains(i.String(), "頁 (100%)") {
			break
		}
	}

	i.CheckAndSendMessages()

	if !forceView {
		if i.previousPost != nil {
			if i.previousPost.ID != i.currentPost.ID {

				err = i.previousPost.Update().SetCurrentViewing(false).Exec(context.Background())
				if err != nil {
					log.Warn(err)
				}

				i.discord.Rest().AddReaction(i.previousPost.ChannelID, i.previousPost.ChannelID, "🈳")
			}
		}

	UpdateLoop:
		for j := 0; j < i.currentSearch.TrackingSeconds; j++ {
			select {
			case <-i.reFetchSignal:
				for len(i.reFetchSignal) > 0 {
					<-i.reFetchSignal
				}
				break UpdateLoop
			default:
				i.SendMessage(event.LiveUpdateCommand)
				c.PrepareWait()
				i.SendRefresh()
				c.Wait(context.Background())
				time.Sleep(1000 * time.Millisecond)
				i.CheckAndSendMessages()
			}
		}
		i.previousPost = i.currentPost
	} else {
		channel, err := i.discord.Rest().GetChannel(i.currentPost.ChannelID)
		if err != nil {
			log.Warn(err)
		}
		if channel.Type() == discord.ChannelTypeGuildPublicThread {
			i.discord.Rest().AddReaction(channel.ID(), channel.ID(), "✅")
		}
	}

	i.Lock()
	i.viewState = ViewStateBacking
	i.Unlock()

	i.Backing(ctx)

}

func (i *Instance) createThreadAndPinnedMessage(messages []string) {
	if i.currentSearch.ForumChannel != 0 {
		channel, err := i.discord.Rest().GetChannel(i.currentSearch.ForumChannel)
		if err != nil {
			log.Warn(err)
			return
		}
		if channel.Type() != discord.ChannelTypeGuildForum {
			log.Fatalf("Thread channel is not a thread channel, actual type is %v", channel.Type())
			return
		}
		//forumChannel := channel.(discord.GuildForumChannel)
		//tags := CheckAndCreateTags(forumChannel)

		post, err := i.discord.Rest().CreatePostInThreadChannel(i.currentSearch.ForumChannel, discord.ThreadChannelPostCreate{
			Name: i.currentPost.Title,

			Message: discord.NewMessageCreateBuilder().SetContentf("- %s\n文章AID: %s\n文章網址: %s", i.currentPost.Title, i.currentPost.Aid, i.currentPost.URL).Build(),
		})
		if err != nil {
			log.Warn(err)
		}

		var messageIds []snowflake.ID

		for j := 0; j < i.currentSearch.PreFillSize; j++ {
			messageCreate := discord.NewMessageCreateBuilder().SetContentf("預留欄位").SetFlags(discord.MessageFlagSuppressNotifications)
			if j < len(messages) {
				messageCreate.SetContent(messages[j])
			}

			msg, err := i.discord.Rest().CreateMessage(post.ID(), messageCreate.Build())
			if err != nil {
				log.Warn(err)
			}
			if msg != nil {
				messageIds = append(messageIds, msg.ID)

				if j == 0 {
					err := i.discord.Rest().PinMessage(post.ID(), msg.ID)
					if err != nil {
						log.Warn(err)
					}
				}
			}
		}

		i.currentPost, err = i.currentPost.Update().SetChannelID(post.ID()).SetContentMessages(messageIds).Save(context.Background())
		if err != nil {
			log.Warn(err)
		}
		return
	}
	if i.currentSearch.TextChannel != 0 {
		channel, err := i.discord.Rest().GetChannel(i.currentSearch.TextChannel)
		if err != nil {
			log.Warn(err)
			return
		}
		if channel.Type() != discord.ChannelTypeGuildText {
			log.Fatalf("Thread channel is not a guild text channel, actual type is %v", channel.Type())
			return
		}

		post, err := i.discord.Rest().CreateThread(i.currentSearch.TextChannel, discord.GuildPublicThreadCreate{
			Name: i.currentPost.Title,
		})
		if err != nil {
			log.Warn(err)
		}

		var messageIds []snowflake.ID

		for j := 0; j < i.currentSearch.PreFillSize; j++ {
			messageCreate := discord.NewMessageCreateBuilder().SetContentf("預留欄位").SetFlags(discord.MessageFlagSuppressNotifications)
			if j < len(messages) {
				messageCreate.SetContent(messages[j])
			}

			msg, err := i.discord.Rest().CreateMessage(post.ID(), messageCreate.Build())
			if err != nil {
				log.Warn(err)
			}
			if msg != nil {
				messageIds = append(messageIds, msg.ID)

				if j == 0 {
					err := i.discord.Rest().PinMessage(post.ID(), msg.ID)
					if err != nil {
						log.Warn(err)
					}
				}
			}
		}

		i.currentPost, err = i.currentPost.Update().SetChannelID(post.ID()).SetContentMessages(messageIds).Save(context.Background())
		if err != nil {
			log.Warn(err)
		}
		return
	}
}

func (i *Instance) CheckAndSendMessages() {
	datas, err := i.db.Message.Query().Where(message.MessageIDEQ(0)).Order(message.ByCreatedAt(sql.OrderAsc())).WithParentPost().All(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	m := make(map[snowflake.ID]ent.Messages)

	for _, data := range datas {
		channelID := data.Edges.ParentPost.ChannelID
		m[channelID] = append(m[channelID], data)
	}

	for channelID, unResizedMessages := range m {
		chunks := SplitMessageByContentLength(unResizedMessages, 1500)
		for _, contents := range chunks {
			messages := lo.Map(contents, func(item *ent.Message, index int) string {
				return item.Content
			})

			messageCreate := discord.NewMessageCreateBuilder().
				SetContent(strings.Join(messages, "\n")).
				SetFlags(discord.MessageFlagSuppressNotifications).
				SetAllowedMentions(&discord.AllowedMentions{}).Build()

			msg, err := i.discord.Rest().CreateMessage(channelID, messageCreate)
			if err != nil {
				log.Warn(err)
				continue
			}

			var ids []int

			for _, content := range contents {
				ids = append(ids, content.ID)
			}

			err = i.db.Message.Update().SetMessageID(msg.ID).Where(message.IDIn(ids...)).Exec(context.Background())
			if err != nil {
				log.Warn(err)
			}
		}
	}
}

func (i *Instance) Backing(ctx context.Context) error {
	log.Infof("Backing %s", i.viewState.String())
	for !strings.Contains(i.String(), "[呼叫器]") && i.viewState == ViewStateBacking {
		i.SendMessage("\x1B[D")
		i.PrepareWait("backing")
		timeout, c := context.WithTimeout(ctx, time.Second*5)
		i.SendMessage("\x0c")
		i.WaitUpdate(timeout, "backing")
		c()
	}
	return nil
}

func (i *Instance) DiscordHook() {
	handler.SyncCommands(i.discord, commands, nil)

	r := handler.New()

	r.ButtonComponent("/publish-content", func(data discord.ButtonInteractionData, e *handler.ComponentEvent) error {
		originalMessage := e.Message

		messageCreateBuilder := discord.NewMessageCreateBuilder().
			SetContent(originalMessage.Content).
			SetEmbeds(originalMessage.Embeds...).
			SetFlags(originalMessage.Flags).
			SetEphemeral(false)

		return e.CreateMessage(messageCreateBuilder.Build())
	})

	r.Autocomplete("/search-user", func(e *handler.AutocompleteEvent) error {
		switch e.Data.Focused().Name {
		case "user":
			userName := e.Data.String("user")
			authorData, err := i.db.Author.Query().Where(author.AuthorIDContainsFold(userName)).Order(author.ByLastSeen(sql.OrderDesc())).All(context.Background())
			if err != nil {
				log.Warn(err)
				return nil
			}
			var choices []discord.AutocompleteChoice

			for _, data := range authorData {
				choices = append(choices, discord.AutocompleteChoiceString{Name: data.AuthorID, Value: data.AuthorID})
			}

			if len(choices) > 25 {
				choices = lo.DropRight(choices, len(choices)-25)
			}

			return e.AutocompleteResult(choices)
		}

		return nil
	})

	r.SlashCommand("/search-message", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		keyword := data.String("keyword")
		messages, err := i.db.Message.Query().Where(message.ContentContainsFold(keyword)).Order(message.ByCreatedAt(sql.OrderDesc())).WithParentPost().All(context.Background())
		if err != nil {
			log.Warn(err)
			return nil
		}
		var contents []string

		var lastChannel snowflake.ID
		for _, message := range messages {
			if lastChannel != message.Edges.ParentPost.ChannelID {
				lastChannel = message.Edges.ParentPost.ChannelID
				contents = append(contents, fmt.Sprintf("於 %s 的留言:", discord.ChannelMention(lastChannel)))
			}
			content := fmt.Sprintf("%s %s", discord.FormattedTimestampMention(message.CreatedAt.Unix(), discord.TimestampStyleShortTime), message.Content)
			contents = append(contents, content)
		}

		//for _, message := range messages {
		//	content := fmt.Sprintf("%s   %s", message.Content)
		//	contents = append(contents, content)
		//}

		chunks := SplitContentByLength(strings.Join(contents, "\n"), 1800)
		msg := fmt.Sprintf("包含關鍵字 `%s` 的留言:\n", keyword) + chunks[0]
		if len(chunks) > 1 {
			msg += "..."
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(msg).SetEphemeral(true).AddActionRow(discord.NewSuccessButton("公開此訊息", "publish-content")).Build())
	})
	r.SlashCommand("/search-user", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		user := data.String("user")
		authorData, err2 := i.db.Author.Query().Where(author.AuthorIDEqualFold(user)).First(context.Background())
		if err2 != nil {
			if ent.IsNotFound(err2) {
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("找不到使用者").SetEphemeral(true).Build())
			}
			log.Warn(err2)
			return nil
		}

		messages, err2 := authorData.QueryMessages().Order(message.ByCreatedAt(sql.OrderDesc())).WithParentPost().All(context.Background())
		if err2 != nil {
			log.Warn(err2)
			return nil
		}

		var contents []string

		var lastChannel snowflake.ID

		for _, message := range messages {
			if lastChannel != message.Edges.ParentPost.ChannelID {
				lastChannel = message.Edges.ParentPost.ChannelID
				contents = append(contents, fmt.Sprintf("於 %s 的留言:", discord.ChannelMention(lastChannel)))
			}
			content := fmt.Sprintf("%s %s", discord.FormattedTimestampMention(message.CreatedAt.Unix(), discord.TimestampStyleShortTime), message.Content)
			contents = append(contents, content)
		}

		chunks := SplitContentByLength(strings.Join(contents, "\n"), 1800)
		msg := fmt.Sprintf("使用者 `%s` 的近期留言:\n", user) + chunks[0]
		if len(chunks) > 1 {
			msg += "..."
		}

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent(msg).SetEphemeral(true).AddActionRow(discord.NewSuccessButton("公開此訊息", "publish-content")).Build())
	})
	// 1140101-晚
	r.SlashCommand("/fetch-title", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		keyword := data.String("search")
		title, ok := data.OptString("title-constraint")
		keywords := strings.Split(keyword, " ")

		var post *ent.PostInfo
		var err error
		if ok {
			post, err = i.db.PostInfo.Query().Where(postinfo.TitleEqualFold(title)).First(context.Background())
			if err == nil {
				err = post.Update().SetSearchKeywords(keywords).SetForceViewExpire(time.Now().Add(30 * time.Minute)).Exec(context.Background())
				if err != nil {
					log.Warn(err)
				}
				i.reFetchSignal <- true
				return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("已強制執行指定標題任務").SetEphemeral(true).Build())
			}
		}

		post, err = i.db.PostInfo.Create().SetTitle(title).SetSearchKeywords(keywords).SetForceViewExpire(time.Now().Add(30 * time.Minute)).Save(context.Background())
		if err != nil {
			log.Warn(err)
		}

		i.reFetchSignal <- true

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("已強制執行指定標題任務").SetEphemeral(true).Build())
	})
	r.SlashCommand("/fetch-aid", func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		aid := data.String("aid")

		current, err := i.db.PostInfo.Query().Where(postinfo.TitleEQ("Force-AID")).Only(context.Background())
		if err != nil {
			if ent.IsNotFound(err) {
				current, err = i.db.PostInfo.Create().SetAid(aid).SetTitle("Force-AID").SetForceViewExpire(time.Now().Add(30 * time.Minute)).Save(context.Background())
				if err != nil {
					log.Warn(err)
				}
				i.reFetchSignal <- true

				return nil
			}
			log.Warn(err)
		}

		err = current.Update().SetAid(aid).SetForceViewExpire(time.Now().Add(30 * time.Minute)).Exec(context.Background())
		if err != nil {
			log.Warn(err)
		}
		i.reFetchSignal <- true

		return e.CreateMessage(discord.NewMessageCreateBuilder().SetContent("已強制執行AID指定任務").SetEphemeral(true).Build())
	})

	i.discord.AddEventListeners(r)
}
func (i *Instance) Start(ctx context.Context) error {
	if err := i.discord.OpenGateway(ctx); err != nil {
		return err
	}
	i.DiscordHook()

	i.RegisterAccountHandler()

	go i.Connect(ctx)

	i.WaitUpdate(ctx, "default")

	i.SendMessage(i.config.PTT.Connection.Username, true)
	i.SendMessage(i.config.PTT.Connection.Password, true)
	return nil
}
