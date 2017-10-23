package main

import (
	//"github.com/sirupsen/logrus"
	"io"
)

const (
	EVENT_TYPE_PAYLOAD = iota
	EVENT_TYPE_CHANGE_WRITE_TARGET
)

type Event struct {
	eventType int

	// PAYLOAD
	payload []byte

	// NEW_WRITE_TARGET
	writeTarget io.Writer

	fileName string
}

func emptyEvent(eventType int) Event {
	return Event{eventType, nil, nil, ""}
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

type Writer struct {
}

func newWriter() Writer {
	return Writer{}
}

func (w *Writer) Start(in chan Event, out chan Event) {
	log.Debug("Start writer")
	var writer io.Writer
	for {
		event, ok := <-in
		if !ok {
			close(out)
			return
		}
		switch event.eventType {
		case EVENT_TYPE_PAYLOAD:
			if writer == nil {
				log.Error("Write target is not initialized (BUG)")
				continue // discard current event
			}
			_, err := writer.Write(event.payload)
			if err != nil {
				log.Panic("Reader goroutine IO failed")
			}
		case EVENT_TYPE_CHANGE_WRITE_TARGET:
			writer = event.writeTarget
		default:
			log.Error("Unknown event type")
		}
	}
}
