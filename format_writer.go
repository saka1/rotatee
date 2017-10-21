package main

import (
	//"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

const (
	EVENT_TYPE_PAYLOAD = iota
	EVENT_TYPE_ROTATE
)

type Event struct {
	eventType int
	payload   []byte
}

func NewEvent(payload []byte) Event {
	return Event{EVENT_TYPE_PAYLOAD, payload}
}

type FormatWriterCtx struct {
	writerGen func(format string) io.Writer
}

func newFormatWriterCtx() FormatWriterCtx {
	return FormatWriterCtx{defaultWriterGen}
}

func startFormatWriter(format string, ch chan Event, ctx FormatWriterCtx) {
	log.Debug("Start rotatee")
	writer := ctx.writerGen(format)
	for {
		event, ok := <-ch
		if !ok {
			return
		}
		switch event.eventType {
		case EVENT_TYPE_PAYLOAD:
			_, err := writer.Write(event.payload)
			if err != nil {
				log.Panic("Reader goroutine IO failed")
			}
		case EVENT_TYPE_ROTATE:
			writer = ctx.writerGen(format)
		default:
			log.Error("Unknown event type")
		}
	}
}

func defaultWriterGen(format string) io.Writer {
	nowFormat := func(format string) string {
		time := time.Now()
		return time.Format(format)
	}
	writer, err := os.OpenFile(nowFormat(format), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Panic("Fail to open destination file")
	}
	return writer
}
