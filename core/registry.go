package core

import "github.com/patyhank/ptd/core/event"

type EventRegistry struct {
	MatchEvents []event.MatchEvent
}

func (r *EventRegistry) RegisterMatchEvent(e event.MatchEvent) {
	r.MatchEvents = append(r.MatchEvents, e)
}

// RegisterDefaultEvents registers the default events
// TODO: Auto generate this
func (r *EventRegistry) RegisterDefaultEvents() {
	r.RegisterMatchEvent(&event.PressAnyKeyEvent{})
	r.RegisterMatchEvent(&event.DuplicateConnectionEvent{})
	r.RegisterMatchEvent(&event.BadLoginNotifyEvent{})
	r.RegisterMatchEvent(&event.MainScreenEvent{})
	r.RegisterMatchEvent(&event.PostInfoScreenEvent{})
	r.RegisterMatchEvent(&event.CommentDataEvent{})
	r.RegisterMatchEvent(&event.ListPostEvent{})
}
