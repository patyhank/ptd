package event

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type PressAnyKeyEvent struct {
}

func (p *PressAnyKeyEvent) Name() string {
	return "任意鍵"
}

func (p *PressAnyKeyEvent) Matched(content string) bool {
	return strings.Contains(content, "按任意鍵繼續]") || strings.Contains(content, "請按任意鍵繼續")
}

type DuplicateConnectionEvent struct {
}

func (p *DuplicateConnectionEvent) Name() string {
	return "重複連線"
}

func (p *DuplicateConnectionEvent) Matched(content string) bool {
	return strings.Contains(content, "您想刪除其他重複登入的連線嗎")
}

type BadLoginNotifyEvent struct {
}

func (p *BadLoginNotifyEvent) Name() string {
	return "登入失敗"
}

func (p *BadLoginNotifyEvent) Matched(content string) bool {
	return strings.Contains(content, "要刪除以上錯誤嘗試的記錄嗎")
}

type MainScreenEvent struct {
	OnlineUsers int64
	Username    string
	CallerState string
}

func (p *MainScreenEvent) Matched(content string) bool {
	return MainScreenRegex.MatchString(content)
}

func (p *MainScreenEvent) FillGroup(content string) {
	matches := MainScreenRegex.FindStringSubmatch(content)
	if matches != nil {
		//p.PTTDate = matches[1]
		//p.PTTTime = matches[2]
		//p.PTTHoliday = matches[3]
		p.OnlineUsers, _ = strconv.ParseInt(matches[1], 10, 64)
		p.Username = matches[2]
		p.CallerState = matches[3]
	}
}

func (p *MainScreenEvent) Name() string {
	return "主畫面"
}

type PostInfoScreenEvent struct {
	PostAID    string
	PostTitle  string
	PostURL    string
	PostValues string
}

func (p *PostInfoScreenEvent) Matched(content string) bool {
	return PostInfoRegex.MatchString(content)
}

func (p *PostInfoScreenEvent) FillGroup(content string) {
	matches := PostInfoRegex.FindStringSubmatch(content)
	if matches != nil {
		p.PostAID = matches[1]
		p.PostTitle = matches[2]
		p.PostURL = matches[3]
		p.PostValues = matches[4]
	}
}

func (p *PostInfoScreenEvent) Name() string {
	return "文章資訊"
}

type CommentDataEvent struct {
	Comments []CommentData
}

func (c *CommentDataEvent) Matched(content string) bool {
	return CommentRegex.MatchString(content)
}

func (c *CommentDataEvent) FillGroup(content string) {
	matches := CommentRegex.FindAllStringSubmatch(content, -1)
	c.Comments = make([]CommentData, len(matches))
	year := time.Now().Year()

	for i, match := range CommentRegex.FindAllStringSubmatch(content, -1) {
		c.Comments[i].Type, _ = ParseCommentType(match[1])

		c.Comments[i].Author = match[2]
		c.Comments[i].Content = match[3]
		if match[4] != "" {
			c.Comments[i].Time, _ = time.Parse("2006/01/02 15:04", fmt.Sprintf("%d/%s", year, match[4])) // fill year to current
		}
	}
}

func (p *CommentDataEvent) Name() string {
	return "留言"
}

type ListPostEvent struct {
	Posts []PostInfo
}

func (l *ListPostEvent) Matched(content string) bool {
	return PostRegex.MatchString(content)
}

func (l *ListPostEvent) FillGroup(content string) {
	matches := PostRegex.FindAllStringSubmatch(content, -1)
	l.Posts = make([]PostInfo, len(matches))

	for i, match := range matches {
		l.Posts[i].Cursor = match[1] == ">"
		l.Posts[i].Number = match[2]
		l.Posts[i].Status = match[3]

		if match[4] == "爆" {
			l.Posts[i].Hot = true
		} else if strings.HasPrefix(match[4], "X") {
			l.Posts[i].Likes, _ = strconv.Atoi(match[4][1:])
			l.Posts[i].Likes *= -10
		} else {
			l.Posts[i].Likes, _ = strconv.Atoi(match[4])
		}

		l.Posts[i].Date = match[5] + "/" + match[6]
		l.Posts[i].Username = match[7]
		l.Posts[i].Title = match[8]
	}
}

func (l *ListPostEvent) Name() string {
	return "文章列表"
}

//type TemplateEvent struct {
//	GenericEvent
//}
//
//func (p *TemplateEvent) Matched(content string) bool {
//	return PostInfoRegex.MatchString(content)
//}
//
//func (p *TemplateEvent) FillGroup(content string) {
//	matches := PostInfoRegex.FindStringSubmatch(content)
//	if matches != nil {
//
//	}
//}
//
//func (p *TemplateEvent) Name() string {
//	return ""
//}
