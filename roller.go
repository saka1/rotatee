package main

import (
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Roller struct {
	historySize int
}

func NewRoller() Roller {
	//TODO
	//return Roller{-1}
	return Roller{3}
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
			log.WithFields(logrus.Fields{"currentFile": currentFile}).Info("Current file closed")
			fallthrough
		case EVENT_TYPE_INIT:
			if historyEnabled {
				lastName := window.slide(event.fileName, func(old string, new string) {
					log.WithFields(logrus.Fields{"old": old, "new": new}).Info("History rotation")
					os.Rename(old, new)
					//TODO error handle
				})
				if lastName != "" {
					log.WithFields(logrus.Fields{"name": lastName}).Info("Remove oldest file at history rotation")
					os.Remove(lastName)
					//TODO error handle
				}
				currentFile = newFile(window.current())
			} else {
				currentFile = newFile(event.fileName)
			}
			log.WithFields(logrus.Fields{"currentFile": currentFile.Name()}).Info("New file opened")
		case EVENT_TYPE_PAYLOAD:
			_, err := currentFile.Write(event.payload)
			if err != nil {
				log.WithFields(logrus.Fields{"err": err}).Panic("Reader goroutine IO failed")
			}
		default:
			log.Warn("Unknown event type")
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

func (hw *historyWindow) current() string {
	//TODO range check
	return evalHistory(hw.names[0], 0)
}

func (hw *historyWindow) last() string {
	if len(hw.names) == 0 {
		return ""
	}
	return evalHistory(hw.names[0], len(hw.names)-1)
}

func (hw *historyWindow) slide(format string, f func(old string, new string)) string {
	if len(hw.names) < hw.limit {
		hw.names = append(hw.names, format)
		return ""
	}
	hw.names = append(hw.names, format)
	for i := 0; i < hw.limit-1; i++ {
		oldName := evalHistory(hw.names[i], len(hw.names)-i-1)
		newName := evalHistory(hw.names[i+1], len(hw.names)-i-2)
		f(oldName, newName)
	}
	last := evalHistory(hw.names[0], hw.limit)
	hw.names = hw.names[1:]
	return last
}

func evalHistory(format string, history int) string {
	r := regexp.MustCompile("([^%])%i") //TODO refactor
	if history == 0 {
		return r.ReplaceAllString(format, "${1}")
	}
	return r.ReplaceAllString(format, "${1}"+strconv.FormatInt(int64(history), 10))
}
