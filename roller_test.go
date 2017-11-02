package main

import (
	"testing"
)

func Test_historyWindow_test(t *testing.T) {
	count := 0
	format := "name%i"
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
	win := newHistoryWindow(3)
	win.slide(format, f)
	t.Log("-------------")
	win.slide(format, f)
	t.Log("-------------")
	win.slide(format, f)
	t.Log("-------------")
	if win.current() != "name" {
		t.Fatalf("unexpected current = %v", win.current())
	}
	if win.last() != "name2" {
		t.Fatalf("unexpected last = %v", win.last())
	}
	expired := win.slide(format, f)
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