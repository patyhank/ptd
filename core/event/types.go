package event

import "time"

type CommentData struct {
	Type    CommentType
	Author  string
	Content string
	Time    time.Time // 精確至分鐘
}

type PostInfo struct {
	Cursor   bool
	Number   string
	Status   string
	Hot      bool
	Likes    int
	Date     string
	Username string
	Title    string
}
