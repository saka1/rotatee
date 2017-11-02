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
	historyEnabled bool
}

func NewRoller() Roller {
	return Roller{-1, false}
}

func NewRollerWithHistory(historySize int) Roller {
	if historySize < -1 {
		panic("Roller: invalid historySize value (BUG)")
	} else if historySize == 0 { // 0 means infinity
		historySize = math.MaxInt32
	}
	return Roller{historySize, historySize != -1}
}

func (roller Roller) Run(in chan Event, out chan Event) {
	var currentFile *os.File = nil
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
			err := currentFile.Close()
			if err != nil {
				log.WithFields(
					logrus.Fields{"err": err.Error(), "name": currentFile.Name(),
					}).Error("Fail to close file when rotation")
			}
			log.WithFields(logrus.Fields{"currentFile": currentFile}).Info("Current file closed")
			fallthrough
		case EVENT_TYPE_INIT:
			if roller.historyEnabled {
				lastName := window.slide(event.fileName, func(old string, new string) {
					log.WithFields(logrus.Fields{"old": old, "new": new}).Info("History rotation")
					err := os.Rename(old, new)
					if err != nil {
						log.WithFields(logrus.Fields{"err": err.Error()}).Error("Fail to rename file when rotation")
					}
				})
				if lastName != "" {
					log.WithFields(logrus.Fields{"name": lastName}).Info("Remove oldest file at history rotation")
					err := os.Remove(lastName)
					if err != nil {
						log.WithFields(logrus.Fields{"name": lastName}).Error("Fail to remove file when rotation")
					}
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

func evalHistory(format string, history int) string {
	r := regexp.MustCompile("([^%])%i") //TODO refactor
	if history == 0 {
		return r.ReplaceAllString(format, "${1}")
	}
	return r.ReplaceAllString(format, "${1}"+strconv.FormatInt(int64(history), 10))
}
