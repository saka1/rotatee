package main

import (
	"github.com/sirupsen/logrus"
	"math"
	"os"
)

type Roller struct {
	historySize int
}

func NewRoller() Roller {
	return Roller{-1}
}

func NewRollerWithHistory(historySize int) Roller {
	if historySize < -1 {
		panic("Roller: ivalid historySize value (BUG)")
	} else if historySize == 0 { // 0 means infinity
		historySize = math.MaxInt32
	}
	return Roller{historySize}
}

func (roller Roller) Run(in chan Event, out chan Event) {
	var currentFile *os.File = nil
	historyEnabled := roller.historySize != -1
	window := newHistoryWindow(roller.historySize)
	for {
		event, ok := <-in
		if !ok {
			if currentFile != nil {
				currentFile.Close()
			}
			close(out)
			return
		}
		switch event.eventType {
		case EVENT_TYPE_CHANGE_WRITE_TARGET:
			currentFile.Close()
			log.WithFields(logrus.Fields{"file": currentFile}).Info("Current file closed")
			fallthrough
		case EVENT_TYPE_INIT:
			f := newFile(event.fileName) //TODO eval fileName( %i to history number )
			log.WithFields(logrus.Fields{"fileName": f.Name()}).Info("New file opened")
			if historyEnabled {
				if name := window.nextPushedOut(); name != "" {
					log.WithFields(logrus.Fields{"name": name}).Info("Remove oldest file at history rotation")
					os.Remove(name)
					//TODO error handle
				}
				window.slide(event.fileName, func(old string, new string) {
					log.WithFields(logrus.Fields{"old": old, "new": new}).Info("History rotation")
					os.Rename(old, new)
					//TODO error handle
				})
			}
			currentFile, event.writeTarget = f, f
			out <- event
		default:
			out <- event
		}
	}
}

func newFile(fileName string) *os.File {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.WithFields(logrus.Fields{"fileName": fileName, "err": err}).Error("Fail to open file")
		return nil
	}
	return file
}

type historyWindow struct {
	limit int
	names []string
}

func newHistoryWindow(limit int) historyWindow {
	return historyWindow{
		limit: limit,
		names: []string{},
	}
}

func (hw *historyWindow) nextPushedOut() string {
	if len(hw.names) < hw.limit {
		return ""
	}
	return hw.names[0]
}

func (hw *historyWindow) slide(name string, f func(old string, new string)) {
	if len(hw.names) < hw.limit {
		hw.names = append(hw.names, name)
		return
	}
	old := hw.names
	new := append(old[1:], name)
	hw.names = new
	for i, o := range old {
		f(o, new[i])
	}
}
