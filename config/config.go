package config

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"strings"
)

// Config is the configuration for the bot
type Config struct {
	Discord Discord   `yaml:"discord"`
	PTT     PTTConfig `yaml:"ptt"`

	Searches []SearchConfig `yaml:"searches"`
}

type Discord struct {
	Token string `yaml:"token"`
}

type EmojiConfig struct {
	UpVote   discord.Emoji `yaml:"up_vote"`
	DownVote discord.Emoji `yaml:"down_vote"`
	Reply    discord.Emoji `yaml:"reply"`
}

type PTTConfig struct {
	Connection ConnectionConfig `yaml:"connection"`
}

type ConnectionConfig struct {
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Host       string `yaml:"host"`
	HostOrigin string `yaml:"host_origin"`
}
type SearchConfig struct {
	Board              string          `json:"board"`
	SearchVariant      []SearchPattern `json:"search_variant"`
	TitleSearchVariant []SearchPattern `json:"title_search_variant"`

	PreFillSize int `yaml:"prefill_size"`

	PostMatchRegex string `json:"post_match_regex"`
	PostTitle      string `json:"post_title"`

	TrackingSeconds int `yaml:"tracking_seconds"`

	ForumChannel snowflake.ID `yaml:"forum_channel,omitempty"`
	TextChannel  snowflake.ID `yaml:"channel,omitempty"`
	Emoji        EmojiConfig  `yaml:"emoji"`
}

type SearchPattern string

func (s SearchPattern) Keys() []string {
	content := string(s)

	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\\r", "\r")

	args := strings.Split(content, "|")

	return args
}
