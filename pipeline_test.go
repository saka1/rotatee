package main

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

type F func(in chan Event, out chan Event)

func (f F) Run(in chan Event, out chan Event) {
	f(in, out)
}

func Test_EventPipe_test(t *testing.T) {
	log.Level = logrus.DebugLevel
	defer func() { log.Level = logrus.InfoLevel }()
	pipeGroup := NewEventPipeGroup()
	pipe := NewEventPipe()
	f := F(func(in chan Event, out chan Event) {
		for {
			input, ok := <-in
			t.Logf("receive input: input = %v, ok = %v", input, ok)
			if !ok {
				close(out)
				return
			}
			if !reflect.DeepEqual(input.payload, []byte("abc")) {
				t.Fatalf("payload mismatch")
			}
		}
	})
	pipe.Add(f)
	pipeGroup.Add(&pipe)
	pipeGroup.Start()
	pipeGroup.Broadcast(NewPayload([]byte("abc")))
	pipeGroup.Stop()
	if !pipeGroup.Closed() {
		t.Fatalf("not closed!!!!")
	}
}
