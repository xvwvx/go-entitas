package entitas

type EventType uint

const (
	EventAdded EventType = iota
	EventUpdated
	EventRemoved

	EventAddedOrRemoved
)
