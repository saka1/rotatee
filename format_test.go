package main

import (
	"testing"
	"time"
)

func Test_format_String(t *testing.T) {
	f := Format("abc")
	if f.String() != "abc" {
		t.Fatal("mismatch String()")
	}
}

func Test_format_evalHistory(t *testing.T) {
	f := Format("abc%i")
	t0 := time.Now()
	// strip %i when history == 0
	if f.evalHistory(t0, 0) != "abc" {
		t.Fatal()
	}
	if f.evalHistory(t0, 1) != "abc1" {
		t.Fatal()
	}
	escaped := Format("abc%%i")
	if escaped.evalHistory(t0, 0) != "abc%i" {
		t.Fatal()
	}
}
