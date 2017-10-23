package main

import (
	"bytes"
	"testing"
)

func Test_writer(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event)
	var buffer bytes.Buffer
	writer := newWriter()
	event := NewWriteTarget()
	event.writeTarget = &buffer
	in <- event
	in <- NewPayload([]byte("hello"))
	close(in)
	writer.Start(in, out)
	if buffer.String() != "hello" {
		t.Fail()
	}
}
