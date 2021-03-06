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
	eos      chan bool
	closed   bool
}

func NewEventPipe() EventPipe {
	return EventPipe{
		stages:   []Stage{},
		channels: []chan Event{},
		eos:      make(chan bool),
		closed:   false,
	}
}

func (e *EventPipe) Add(stage Stage) {
	e.stages = append(e.stages, stage)
	e.channels = append(e.channels, make(chan Event))
}

func (e *EventPipe) Start() {
	// create "tail" channel
	// all messages should NOT be handled by this process
	tailCh := make(chan Event)
	go func() {
		for {
			event, ok := <-tailCh
			if !ok {
				e.eos <- true
				log.Debug("tail of pipeline closed")
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

func (e *EventPipe) InputCh() chan Event {
	return e.channels[0]
}

func (e *EventPipe) Stop() {
	close(e.InputCh())
	e.closed = <-e.eos
}

func (e *EventPipe) Closed() bool {
	return e.closed
}

type EventPipeGroup struct {
	pipes []*EventPipe
}

func NewEventPipeGroup() EventPipeGroup {
	return EventPipeGroup{}
}

func (e *EventPipeGroup) Add(pipe *EventPipe) {
	e.pipes = append(e.pipes, pipe)
}

func (e *EventPipeGroup) Start() {
	for _, p := range e.pipes {
		p.Start()
	}
}

func (e *EventPipeGroup) Broadcast(ev Event) {
	for _, p := range e.pipes {
		p.InputCh() <- ev
	}
}

func (e *EventPipeGroup) Stop() {
	for _, p := range e.pipes {
		p.Stop()
	}
}

func (e *EventPipeGroup) Closed() bool {
	for _, p := range e.pipes {
		if !p.Closed() {
			return false
		}
	}
	return true
}
