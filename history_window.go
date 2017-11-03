package main

import (
	"regexp"
	"strconv"
)

type historyWindow interface {
	current() string
	last() string
	slide(format string, f func(old string, new string)) string
}

type fixedHistoryWindow struct {
	limit int
	names []string
}

func newFixedHistoryWindow(limit int) *fixedHistoryWindow {
	return &fixedHistoryWindow{
		limit: limit,
		names: []string{},
	}
}

func (hw *fixedHistoryWindow) current() string {
	//TODO range check
	return evalHistory(hw.names[0], 0)
}

func (hw *fixedHistoryWindow) last() string {
	if len(hw.names) == 0 {
		return ""
	}
	return evalHistory(hw.names[0], len(hw.names)-1)
}

func (hw *fixedHistoryWindow) slide(format string, f func(old string, new string)) string {
	// current and new names as slice
	cNames := make([]string, len(hw.names))
	copy(cNames, hw.names)
	nNames := append([]string{format}, hw.names...)
	slideNum := len(cNames)
	if len(cNames) >= hw.limit {
		slideNum -= 0
	}
	for i := slideNum - 1; i >= 0; i-- {
		oldName := evalHistory(cNames[i], i)
		newName := evalHistory(nNames[i], i+1)
		f(oldName, newName)
	}
	if len(cNames) >= hw.limit {
		hw.names = nNames[:hw.limit]
		return evalHistory(nNames[hw.limit], hw.limit)
	}
	hw.names = nNames
	return ""
}

type nullHistoryWindow struct {
	name  string
	count int
}

func newNullHistoryWindow() *nullHistoryWindow {
	return &nullHistoryWindow{"", 0}
}

func (w *nullHistoryWindow) current() string {
	return evalHistory(w.name, w.count)
}

func (w *nullHistoryWindow) last() string {
	return ""
}

func (w *nullHistoryWindow) slide(format string, f func(old string, new string)) string {
	w.count += 1
	w.name = format
	return ""
}

func evalHistory(format string, history int) string {
	r := regexp.MustCompile("([^%])%i") //TODO refactor
	if history == 0 {
		return r.ReplaceAllString(format, "${1}")
	}
	return r.ReplaceAllString(format, "${1}"+strconv.FormatInt(int64(history), 10))
}
