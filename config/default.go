package config

import (
	"github.com/disgoorg/disgo/discord"
)

var DefaultConfig = Config{
	Discord: Discord{
		Token: "",
	},
	PTT: PTTConfig{
		Connection: ConnectionConfig{
			Username:   "user",
			Password:   "pass",
			Host:       "wss://ws.ptt.cc/bbs",
			HostOrigin: "https://term.ptt.cc",
		},
	},
	Searches: []SearchConfig{
		{
			Board:              "c_chat",
			SearchVariant:      []SearchPattern{"/holo\\r/直播\\r", "/holo\\r/間直播\\r", "/holo\\r/直播單\\r"}, // 不建議超過3避免無法搜尋文章日期等
			TitleSearchVariant: []SearchPattern{"/holo\\r/直播\\r", "/holo\\r/間直播\\r", "/holo\\r/直播單\\r"},
			TrackingSeconds:    600,
			PostMatchRegex:     "\\[\\S+?] \\[Vtub] {1,2}Hololive ([晚日])間直播單(?:（| \\()(\\d+)(?:）|\\))",
			PostTitle:          "%[2]s-%[1]s間直播串",

			ForumChannel: 0,
			TextChannel:  0,
			Emoji: EmojiConfig{
				UpVote: discord.Emoji{
					Name: "push",
					ID:   1033734367937310720,
				},
				DownVote: discord.Emoji{
					Name: "boo",
					ID:   1033734366330892428,
				},
				Reply: discord.Emoji{
					Name: "addon",
					ID:   1033734364950953994,
				},
			},
		},
	},
}
