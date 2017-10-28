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
		if count == 1 && !(old == "hoge-" && new == "hoge-1") {
			t.Fatal()
		}
	}
	win := newHistoryWindow(2)
	win.slide(format, f)
	win.slide(format, f)
	if win.current() != "hoge-" { //TODO fix
		t.Fatalf("unexpected current = %v", win.current())
	}
	if win.last() != "hoge-1" {
		t.Fatalf("unexpected last = %v", win.last())
	}
	win.slide(format, f)
	if count != 1 {
		t.Fatalf("unexpected number of calls = %v", count)
	}
	if len(win.names) != 2 {
		t.Fatalf("unexpected len of names = %v", len(win.names))
	}
}
