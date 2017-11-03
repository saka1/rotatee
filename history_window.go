package main

import "time"

type historyWindow interface {
	current() string
	last() string
	slide(format Format, t time.Time, f func(old string, new string)) string
}

type fixedHistoryWindow struct {
	limit int
	names []Format
	times []time.Time
}

func newFixedHistoryWindow(limit int) *fixedHistoryWindow {
	return &fixedHistoryWindow{
		limit: limit,
		names: []Format{},
		times: []time.Time{},
	}
}

func (hw *fixedHistoryWindow) current() string {
	return hw.names[0].evalHistory(hw.times[0], 0)
}

func (hw *fixedHistoryWindow) last() string {
	if len(hw.names) == 0 {
		return ""
	}
	return hw.names[0].evalHistory(hw.times[0], len(hw.names)-1)
}

func (hw *fixedHistoryWindow) slide(format Format, t time.Time, f func(old string, new string)) string {
	// current and new names as slice
	cNames, cTimes := make([]Format, len(hw.names)), make([]time.Time, len(hw.names))
	copy(cNames, hw.names)
	copy(cTimes, hw.times)
	nNames, nTimes := append([]Format{format}, hw.names...), append([]time.Time{t}, hw.times...)
	slideNum := len(cNames)
	if len(cNames) >= hw.limit {
		slideNum -= 0
	}
	for i := slideNum - 1; i >= 0; i-- {
		oldName := cNames[i].evalHistory(cTimes[i], i)
		newName := nNames[i].evalHistory(nTimes[i], i+1)
		f(oldName, newName)
	}
	if len(cNames) >= hw.limit {
		hw.names, hw.times = nNames[:hw.limit], nTimes[:hw.limit]
		return nNames[hw.limit].evalHistory(nTimes[hw.limit], hw.limit)
	}
	hw.names, hw.times = nNames, nTimes
	return ""
}

type nullHistoryWindow struct {
	name Format
}

func newNullHistoryWindow() *nullHistoryWindow {
	return &nullHistoryWindow{""}
}

func (w *nullHistoryWindow) current() string {
	return w.name.String()
}

func (w *nullHistoryWindow) last() string {
	return ""
}

func (w *nullHistoryWindow) slide(format Format, t time.Time, f func(old string, new string)) string {
	w.name = format
	return ""
}
