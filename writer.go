package main

import (
	//"github.com/sirupsen/logrus"
	"io"
)

type Writer struct {
}

func NewWriter() Writer {
	return Writer{}
}

func (w Writer) Run(in chan Event, out chan Event) {
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
