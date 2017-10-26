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
			if currentFile != nil { // MEMO: currentFile is nil when the first time
				currentFile.Close()
			}
			currentFile = f
			event.writeTarget = f
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