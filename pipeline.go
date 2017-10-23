package main

import (
	"github.com/sirupsen/logrus"
)

// This type represents a stage of pipeline.
// All stages is called async(on goroutine)
type Stage interface {
	Run(in chan Event, out chan Event)
}

type EventPipe struct {
	stages   []Stage
	channels []chan Event
}

func NewEventPipe() EventPipe {
	return EventPipe{
		stages:   []Stage{},
		channels: []chan Event{},
	}
}

func (e *EventPipe) Add(stage Stage) {
	e.stages = append(e.stages, stage)
	e.channels = append(e.channels, make(chan Event))
}

func (e *EventPipe) Start(ev Event) {
	// create "tail" channel
	// all messages received by the channel will be dropped
	tailCh := make(chan Event)
	// drop message prosess(used as tail process)
	go func() {
		for {
			event, ok := <-tailCh
			if !ok {
				//TODO detect shutdown here
				return
			}
			log.WithFields(logrus.Fields{
				"eventType": event.eventType,
			}).Warn("Unhandled message found at the end of pipeline")
		}
	}()
	e.channels = append(e.channels, tailCh)
	for i, stage := range e.stages {
		go stage.Run(e.channels[i], e.channels[i+1])
	}
}
