package main

import (
	"io"
	"time"
)

const (
	EVENT_TYPE_PAYLOAD = iota
	EVENT_TYPE_CHANGE_WRITE_TARGET
)

type Event struct {
	eventType   int
	payload     []byte
	writeTarget io.Writer
	fileName    string
	timestamp   time.Time
}

func emptyEvent(eventType int) Event {
	return Event{eventType, nil, nil, "", time.Now()}
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
