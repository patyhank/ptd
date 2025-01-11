package common

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/patyhank/ptd/ent"
	"sync"
)

type Instance struct {
	discord      bot.Client
	currentPost  *ent.PostInfo // 目前瀏覽的貼文
	previousPost *ent.PostInfo // 上一個瀏覽的貼文，瀏覽下一篇時，於Discord新增已完成標籤
	db           *ent.Client

	viewState ViewState
	sync.Mutex
}

func NewInstance() {

}
