package main

import (
	"github.com/sirupsen/logrus"
	"os"
	"math"
)

type Roller struct {
	window historyWindow
	setting RotateeSetting
}

func NewRoller(format Format, setting RotateeSetting) Roller {
	historySize := setting.historySize
	if format.HasHistoryNumberSpec() {
		window := newFixedHistoryWindow(historySize)
		if historySize == 0 {
			log.Warn("Unlimited history size may cause performance impact. " +
				"Consider to use `--history` option")
			window = newFixedHistoryWindow(math.MaxInt32)
		}
		return Roller{window: window}
	}
	if historySize == 0 {
		return Roller{newNullHistoryWindow(), setting}
	}
	if historySize > 32 { //TODO rethink this threshold
		log.Warn("Large size of history may cause performance impact. " +
			"Consider to use more smaller size")
	}
	return Roller{newFixedHistoryWindow(historySize), setting}
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
			lastName := roller.window.slide(event.format, event.timestamp, func(old string, new string) {
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
			}
			currentFile = newFile(roller.window.current(), roller.setting.appendMode)
			log.WithFields(logrus.Fields{"currentFile": currentFile, "name": currentFile.Name()}).Info("New file opened")
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

func newFile(fileName string, appendMode bool) *os.File {
	flag := os.O_RDWR|os.O_CREATE
	if appendMode {
		flag |= os.O_APPEND
	}
	file, err := os.OpenFile(fileName, flag, 0644)
	if err != nil {
		log.WithFields(logrus.Fields{"fileName": fileName, "err": err}).Error("Fail to open file")
		return nil
	}
	return file
}
