package event

import (
	"regexp"
	"strings"
)

// LiveUpdateCommand is the command to quick update the live post
const LiveUpdateCommand = "\x1B[D\x1B[C\x1B[4~"

// PostRegex [Message Cursor ID Status PushCount]
var PostRegex = regexp.MustCompile("([> ])(\\d+) ([+Mm -!]) ?(X?\\d+|[爆 ]) ?(\\d{1,2})/ ?(\\d{1,2}) (\\w+)\\s+([\\S\\s]*?)\\s+\\n")

var CommentRegex = regexp.MustCompile("([→推噓])  ?(\\w+): ([\\s\\S]+?)(?: +?| {0})(?:\\S+\\n|(\\d+/\\d+ \\d+:\\d+))")

var MainScreenRegex = regexp.MustCompile("\\[(\\d+/\\d+) (\\S+) (\\d+:\\d+)]\\s+(\\S+)?\\s+線上(\\d+)人, 我是(\\w+)\\s+\\[呼叫器](\\S+)")

var PostInfoRegex = regexp.MustCompile("文章代碼\\(AID\\): #(\\w+) ([\\S\\s]+?)\\s+[│|][\\S\\s]+文章網址: (https://\\S+)[\\s\\S]+ 這一篇文章值 (\\d+) Ptt幣")

func containsAllConditions(content string, cond []string) bool {
	for _, s := range cond {
		if !strings.Contains(content, s) {
			return false
		}
	}
	return true
}

type CommentType int32

const (
	_ CommentType = iota
	CommentTypeUpVote
	CommentTypeDownVote
	CommentTypeReply
)

func (c CommentType) String() string {
	switch c {
	case CommentTypeUpVote:
		return "推"
	case CommentTypeDownVote:
		return "噓"
	case CommentTypeReply:
		return "→"
	}
	return ""
}

func ParseCommentType(content string) (CommentType, bool) {
	switch content {
	case "推":
		return CommentTypeUpVote, true
	case "噓":
		return CommentTypeDownVote, true
	case "→":
		return CommentTypeReply, true
	}
	return 0, false
}

func MustParseCommentType(content string) CommentType {
	t, ok := ParseCommentType(content)
	if !ok {
		panic("invalid comment type")
	}
	return t
}
