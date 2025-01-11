package main

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/patyhank/ptd/config"
	"github.com/patyhank/ptd/ent"
	"sync"
)

type State int32

const (
	StateBacking State = iota
	StateSearching
	StateReadyViewing
	StateViewing
)

const LiveUpdateCommand = "\x1B[D\x1B[C\x1B[4~"

var currentSearchIndex = 0 // index of searchWordVariant
var prefillSize = 15       // number of messages to prefill

// state of the bot
var cfg config.Config
var s bot.Client
var current *ent.PostInfo
var last *ent.PostInfo
var state State = StateBacking
var stateLock = sync.Mutex{}
var dbClient *ent.Client

// types definition
