package common

import "regexp"

const LiveUpdateCommand = "\x1B[D\x1B[C\x1B[4~"

type ViewState int32

const (
	ViewStateBacking ViewState = iota
	ViewStateSearching
	ViewStateReadyViewing
	ViewStateViewing
)

// PostRegex [Message Cursor ID Status PushCount]
var PostRegex = regexp.MustCompile("([> ])(\\d+) ([+Mm -!]) ?(X?\\d+|[爆 ]) ?(\\d{1,2})/ ?(\\d{1,2}) (\\w+)\\s+([\\S\\s]*?)\\s+\\n")

// CommentRegex [Message StrType Author Content Time]
var CommentRegex = regexp.MustCompile("(→ |→|推|噓) (\\w+): ([\\s\\S]+?)(?: +?| {0})(?:\\S+\\n|(\\d+/\\d+ \\d+:\\d+))")
