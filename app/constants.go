package app

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
)

type ViewState int32

const (
	ViewStateInit ViewState = iota
	ViewStateBacking
	ViewStateSearching
	ViewStateReadyViewing
	ViewStateViewing
)

func (v ViewState) String() string {
	return [...]string{"Init", "Backing", "Searching", "ReadyViewing", "Viewing"}[v]
}

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name: "fetch-aid",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "搜尋文章",
		},
		Description: "fetch aid for post",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "指定目前使用的AID",
		},
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				Name:        "aid",
				Description: "post aid",
				Required:    true,
			},
		},
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
	},
	discord.SlashCommandCreate{
		Name: "fetch-title",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "搜尋文章標題",
		},
		Description: "fetch post for title",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "",
		},
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				Name:        "search",
				Description: "使用格式 a|b 來搜尋文章關鍵詞a及b",
				Required:    true,
			},
			&discord.ApplicationCommandOptionString{
				Name:        "title-constraint",
				Description: "標題限制",
				Required:    false,
			},
		},
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
	},
	discord.SlashCommandCreate{
		Name: "search-user",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "搜尋用戶",
		},
		Description: "fetch comment for user",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "搜尋用戶的留言",
		},
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				Name:         "user",
				Description:  "post user",
				Required:     true,
				Autocomplete: true,
			},
		},
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
	},
	discord.SlashCommandCreate{
		Name: "search-message",
		NameLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "搜尋留言",
		},
		Description: "fetch comment for keyword",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleChineseTW: "搜尋包含此關鍵字的留言",
		},
		Options: []discord.ApplicationCommandOption{
			&discord.ApplicationCommandOptionString{
				DescriptionLocalizations: map[discord.Locale]string{
					discord.LocaleChineseTW: "留言關鍵字",
				},
				NameLocalizations: map[discord.Locale]string{
					discord.LocaleChineseTW: "關鍵字",
				},
				Name:        "keyword",
				Description: "post keyword",
				Required:    true,
			},
		},
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
	},
}
