package main

import (
	"testing"
)

func Test_historyWindow_test(t *testing.T) {
	count := 0
	format := "hoge-%i"
	f := func(old string, new string) {
		t.Logf("old, new = %v, %v", old, new)
		count += 1
		if count == 1 && !(old == "hoge-3" && new == "hoge-2") {
			t.Fatal()
		}
		if count == 2 && !(old == "hoge-2" && new == "foo-1") {
			t.Fatal()
		}
	}
	win := newHistoryWindow(3)
	win.slide(format, f)
	win.slide(format, f)
	win.slide("foo-%i", f)
	if win.current() != "hoge-" {
		t.Fatalf("unexpected current = %v", win.current())
	}
	if win.last() != "hoge-2" {
		t.Fatalf("unexpected last = %v", win.last())
	}
	win.slide(format, f)
	if count != 2 {
		t.Fatalf("unexpected number of calls = %v", count)
	}
	if len(win.names) != 3 {
		t.Fatalf("unexpected len of names = %v", len(win.names))
	}
}
