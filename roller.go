package main

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Roller struct {
}

func NewRoller() Roller {
	return Roller{}
}

func (roller Roller) Run(in chan Event, out chan Event) {
	var currentFile *os.File = nil
	for {
		event, ok := <-in
		if !ok {
			closeCurrentFile(currentFile)
			close(out)
			return
		}
		switch event.eventType {
		case EVENT_TYPE_CHANGE_WRITE_TARGET:
			f := newFile(event.fileName)
			if f == nil {
				log.Error("Discard internal event bacause of failure to create new file")
				continue
			}
			log.WithFields(logrus.Fields{"fileName": f.Name()}).Info("New file opened")
			closeCurrentFile(currentFile)
			currentFile = f
			event.writeTarget = f
			out <- event
		default:
			out <- event
		}
	}
}

func closeCurrentFile(file *os.File) {
	if file != nil { // MEMO: currentFile is nil when the first time
		file.Close()
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
	names []string
}

func newHistoryWindow() historyWindow {
	return historyWindow{names: []string{}}
}

func (hw *historyWindow) add(name string) {
	hw.names = append(hw.names, name)
}
