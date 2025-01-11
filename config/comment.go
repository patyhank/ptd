package config

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
