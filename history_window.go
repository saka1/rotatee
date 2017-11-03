package main

type historyWindow interface {
	current() string
	last() string
	slide(format Format, f func(old string, new string)) string
}

type fixedHistoryWindow struct {
	limit int
	names []Format
}

func newFixedHistoryWindow(limit int) *fixedHistoryWindow {
	return &fixedHistoryWindow{
		limit: limit,
		names: []Format{},
	}
}

func (hw *fixedHistoryWindow) current() string {
	return hw.names[0].evalHistory(0)
}

func (hw *fixedHistoryWindow) last() string {
	if len(hw.names) == 0 {
		return ""
	}
	return hw.names[0].evalHistory(len(hw.names)-1)
}

func (hw *fixedHistoryWindow) slide(format Format, f func(old string, new string)) string {
	// current and new names as slice
	cNames := make([]Format, len(hw.names))
	copy(cNames, hw.names)
	nNames := append([]Format{format}, hw.names...)
	slideNum := len(cNames)
	if len(cNames) >= hw.limit {
		slideNum -= 0
	}
	for i := slideNum - 1; i >= 0; i-- {
		oldName := cNames[i].evalHistory(i)
		newName := nNames[i].evalHistory(i+1)
		f(oldName, newName)
	}
	if len(cNames) >= hw.limit {
		hw.names = nNames[:hw.limit]
		return nNames[hw.limit].evalHistory(hw.limit)
	}
	hw.names = nNames
	return ""
}

type nullHistoryWindow struct {
	name  Format
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

func (w *nullHistoryWindow) slide(format Format, f func(old string, new string)) string {
	w.name = format
	return ""
}

