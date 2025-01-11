package app

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/patyhank/ptd/app/ent"
	"github.com/samber/lo"
	"strings"
)

// SplitMessageByContentLength 將訊息依照字數分割，避免訊息過長無法發送至Discord
func SplitMessageByContentLength(collection ent.Messages, wordCount int) (chunks []ent.Messages) {
	var chunk []*ent.Message
	var count int
	for _, data := range collection {
		count += len(data.Content)
		if count > wordCount {
			chunks = append(chunks, chunk)
			chunk = []*ent.Message{}
			count = 0
		}
		chunk = append(chunk, data)
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return
}

// SplitContentByLength
// 將字串依照指定的字數分割，並提前於換行處分割
// 這讓Discord訊息更好閱讀，並且不會截斷連結
func SplitContentByLength(input string, wordCount int) (chunks []string) {
	var chunk string
	var count int
	collection := strings.Split(input, "\n")
	for _, data := range collection {
		count += len([]rune(data))
		if count >= wordCount {
			chunks = append(chunks, chunk)
			chunk = ""
			count = 0
		}
		chunk += data + "\n"
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return
}

// Deprecated: Unwrap 用於快速檢查錯誤，如果錯誤不為nil，則會panic，謹慎使用
func Unwrap[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func passingStringsToFmt(format string, strArgs []string) string {
	args := lo.Map(strArgs, func(item string, index int) any {
		return any(item)
	})
	return fmt.Sprintf(format, args...)
}

func emojiName(emoji discord.Emoji) string {
	if emoji.ID == 0 {
		return emoji.Name
	}
	return emoji.Mention()
}
