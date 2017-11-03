package main

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

var log = logrus.New()

type Rotatee struct {
	setting RotateeSetting
}

func NewRotatee(setting RotateeSetting) *Rotatee {
	return &Rotatee{setting: setting}
}

type RotateeSetting struct {
	args        []string
	verbose     bool
	maxFileSize int64
	historySize int
}

func setupEventPipe(setting RotateeSetting) EventPipeGroup {
	pipeGroup := NewEventPipeGroup()
	for _, arg := range setting.args {
		pipe := NewEventPipe()
		if setting.maxFileSize != 0 {
			pipe.Add(NewScaler(setting.maxFileSize))
		}
		pipe.Add(NewTimer(DetectSeries(arg, time.Now())))
		pipe.Add(NewFormatEval(arg))
		pipe.Add(NewRoller(Format(arg), setting.historySize))
		pipeGroup.Add(&pipe)
	}
	return pipeGroup
}

func teeLoop(pipeGroup *EventPipeGroup) {
	reader := os.Stdin
	readBuf := make([]byte, 1024)
	for {
		length, err := reader.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				pipeGroup.Stop()
				os.Exit(0)
			}
			log.Panic("Writer goroutine IO failed")
		}
		// do copy because 'content' is shared among goroutine(s)
		content := make([]byte, length)
		copy(content, readBuf[:length])
		pipeGroup.Broadcast(NewPayload(content))
		// Also, write to stdout
		os.Stdout.Write(content)
	}
}

func (r *Rotatee) Start() {
	log.WithFields(logrus.Fields{"Rotatee": r}).Debug("Start rotatee")
	pipeGroup := setupEventPipe(r.setting)
	pipeGroup.Start()
	pipeGroup.Broadcast(NewInit())
	teeLoop(&pipeGroup)
}
