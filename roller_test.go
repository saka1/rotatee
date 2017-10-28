package main

import (
	"testing"
)

func Test_historyWindow_test(t *testing.T) {
	win := newHistoryWindow(2)
	count := 0
	f := func(old string, new string) {
		count += 1
		t.Logf("old, new = %v, %v", old, new)
		if old == "a" && new == "b" && count == 1 {
			// OK
		} else if old == "b" && new == "c" && count == 2 {
			// OK
		} else {
			t.Fatal("invald old-new pair")
		}
	}
	if win.nextPushedOut() != "" {
		t.Fatal("nextPushedOut is broken")
	}
	// expected window: "a", "b"
	win.slide("a", f)
	win.slide("b", f)
	if win.nextPushedOut() != "a" {
		t.Fatal("nextPushedOut is broken")
	}
	// expected window: "b", "c"
	// expected old-new pair:
	//   "a" -> "b"
	//   "b" -> "c"
	win.slide("c", f)
	if count != 2 {
		t.Fatalf("unexpected call count=%d", count)
	}
	if win.names[0] != "b" || win.names[1] != "c" {
		t.Fatal("unexpected names order")
	}
}
