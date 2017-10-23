package main

import (
	"testing"
)

func Test_scaler_nonInterleave(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event, 32)
	scaler := newScaler(5)
	in <- NewPayload([]byte("hello"))
	close(in)
	scaler.Start(in, out)
	if ev := <-out; string(ev.payload) != "hello" {
		t.Fail()
	}
	if _, ok := <-out; ok {
		t.Fail()
	}
}

func Test_scaler_interleave(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event, 32)
	scaler := newScaler(4)
	in <- NewPayload([]byte("hello"))
	close(in)
	scaler.Start(in, out)
	if ev := <-out; string(ev.payload) != "hell" {
		t.Fail()
	}
	if ev := <-out; ev.eventType != EVENT_TYPE_CHANGE_WRITE_TARGET {
		t.Fail()
	}
	if ev := <-out; string(ev.payload) != "o" {
		t.Fail()
	}
	if _, ok := <-out; ok {
		t.Fail()
	}
}
