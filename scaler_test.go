package main

import (
	"testing"
)

func Test_scaler_nonInterleave(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event, 32)
	scaler := NewScaler(5)
	in <- NewPayload([]byte("hello"))
	close(in)
	scaler.Run(in, out)
	if ev := <-out; string(ev.payload) != "hello" {
		t.Fatal()
	}
	if _, ok := <-out; ok {
		t.Fatal()
	}
}

func Test_scaler_interleave(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event, 32)
	scaler := NewScaler(4)
	in <- NewPayload([]byte("hello"))
	close(in)
	scaler.Run(in, out)
	if ev := <-out; string(ev.payload) != "hell" {
		t.Fatal()
	}
	if ev := <-out; ev.eventType != EVENT_TYPE_CHANGE_WRITE_TARGET {
		t.Fatal()
	}
	if ev := <-out; string(ev.payload) != "o" {
		t.Fatal()
	}
	if _, ok := <-out; ok {
		t.Fatal()
	}
}

func Test_scaler_large(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event, 32)
	scaler := NewScaler(2)
	in <- NewPayload([]byte("hello"))
	close(in)
	scaler.Run(in, out)
	if ev := <-out; string(ev.payload) != "he" {
		t.Fatalf("invalid receive: %v", string(ev.payload))
	}
	if ev := <-out; ev.eventType != EVENT_TYPE_CHANGE_WRITE_TARGET {
		t.Fatal()
	}
	if ev := <-out; string(ev.payload) != "ll" {
		t.Fatalf("invalid receive: %v", string(ev.payload))
	}
	if ev := <-out; ev.eventType != EVENT_TYPE_CHANGE_WRITE_TARGET {
		t.Fatal()
	}
	if ev := <-out; string(ev.payload) != "o" {
		t.Fatal()
	}
	if _, ok := <-out; ok {
		t.Fatal()
	}
}

func Test_scaler_smallMsg(t *testing.T) {
	in, out := make(chan Event, 32), make(chan Event, 32)
	scaler := NewScaler(100)
	for i := 0; i < 10; i++ {
		in <- NewPayload([]byte("hello")) // 5byte * 10msg
	}
	close(in)
	scaler.Run(in, out)
	if scaler.count != 50 {
		t.Fatalf("invalid counter = %v", scaler.count)
	}
}
