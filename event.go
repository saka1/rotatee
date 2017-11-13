package main

import (
	"time"
)

type EventType int

const (
	EventTypePayload EventType = iota
	EventTypeChangeWriteTarget
	EventTypeInit
)

type Event struct {
	eventType EventType
	payload   []byte
	format    Format
	timestamp time.Time
}

func emptyEvent(eventType EventType) Event {
	return Event{eventType, nil, Format(""), time.Now()}
}

func NewPayload(payload []byte) Event {
	ev := emptyEvent(EventTypePayload)
	ev.payload = payload
	return ev
}

func NewWriteTarget() Event {
	ev := emptyEvent(EventTypeChangeWriteTarget)
	return ev
}

func NewInit() Event {
	return emptyEvent(EventTypeInit)
}
