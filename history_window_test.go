package main

import (
	"testing"
	"time"
)

func Test_historyWindow_test(t *testing.T) {
	count := 0
	format := Format("name%i")
	t0 := time.Now()
	f := func(old string, new string) {
		t.Logf("old, new = %v, %v", old, new)
		count += 1
		if count == 1 && !(old == "name" && new == "name1") {
			t.Fatal()
		}
		if count == 2 && !(old == "name1" && new == "name2") {
			t.Fatal()
		}
		if count == 3 && !(old == "name" && new == "name1") {
			t.Fatal()
		}
	}
	win := newFixedHistoryWindow(3)
	win.slide(format, t0, f)
	t.Log("-------------")
	win.slide(format, t0, f)
	t.Log("-------------")
	win.slide(format, t0, f)
	t.Log("-------------")
	if win.current() != "name" {
		t.Fatalf("unexpected current = %v", win.current())
	}
	if win.last() != "name2" {
		t.Fatalf("unexpected last = %v", win.last())
	}
	expired := win.slide(format, t0, f)
	if count != (1 + 2 + 3) {
		t.Fatalf("unexpected number of calls = %v", count)
	}
	if len(win.names) != 3 {
		t.Fatalf("unexpected len of names = %v, names = %v", len(win.names), win.names)
	}
	if expired != "name3" {
		t.Fatalf("unexpected return from slide() = %v", expired)
	}
}
