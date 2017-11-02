package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
)

type Roller struct {
	window historyWindow
}

func NewRoller() Roller {
	return Roller{newNullHistoryWindow()}
}

func NewRollerWithHistory(historySize int) Roller {
	window := newFixedHistoryWindow(historySize)
	return Roller{window}
}

func (roller Roller) Run(in chan Event, out chan Event) {
	var currentFile *os.File = nil

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
			lastName := roller.window.slide(event.fileName, func(old string, new string) {
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
			currentFile = newFile(roller.window.current())
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
	name string
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
