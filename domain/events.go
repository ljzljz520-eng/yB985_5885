package domain

import "time"

type Event struct {
	ID, AggregateID, Type, Actor, Payload string
	At                                    time.Time
}

func NewEvent(id, aggregate, typ, actor, payload string, at time.Time) Event {
	return Event{ID: id, AggregateID: aggregate, Type: typ, Actor: actor, Payload: payload, At: at}
}
func EventLabel(e Event) string {
	if e.Actor == "" {
		return e.Type
	}
	return e.Type + " by " + e.Actor
}
