package event

type RegexMatchEvent interface {
	MatchEvent
	FillGroup(content string)
}

type MatchEvent interface {
	Event
	Matched(content string) bool
}

type Event interface {
	Name() string
}
