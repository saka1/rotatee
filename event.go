package main

import (
	"time"
)

const (
	EVENT_TYPE_PAYLOAD             = iota
	EVENT_TYPE_CHANGE_WRITE_TARGET
	EVENT_TYPE_INIT
)

type Event struct {
	eventType int
	payload   []byte
	format    Format
	timestamp time.Time
}

func emptyEvent(eventType int) Event {
	return Event{eventType, nil, Format(""), time.Now()}
}

func NewPayload(payload []byte) Event {
	ev := emptyEvent(EVENT_TYPE_PAYLOAD)
	ev.payload = payload
	return ev
}

func NewWriteTarget() Event {
	ev := emptyEvent(EVENT_TYPE_CHANGE_WRITE_TARGET)
	return ev
}

func NewInit() Event {
	return emptyEvent(EVENT_TYPE_INIT)
}
