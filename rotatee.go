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
	args    []string
	verbose bool
}

func setupEventPipe(setting RotateeSetting) *EventPipeGroup {
	pipeGroup := NewEventPipeGroup()
	for _, arg := range setting.args {
		pipe := NewEventPipe()
		pipe.Add(NewTimer(DetectSeries(arg, time.Now())))
		pipe.Add(NewFormatEval(arg))
		pipe.Add(NewRoller())
		pipe.Add(NewWriter())
		pipeGroup.Add(pipe)
	}
	return pipeGroup
}

func teeLoop(pipeGroup *EventPipeGroup) {
	reader := os.Stdin
	readBuf := make([]byte, 1024)
	for {
		len, err := reader.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				pipeGroup.Stop()
				os.Exit(0)
			}
			log.Panic("Writer goroutine IO failed")
		}
		// do copy because 'content' is shared among goroutine(s)
		content := make([]byte, len)
		copy(content, readBuf[:len])
		pipeGroup.Broadcast(NewPayload(content))
		// Also, write to stdout
		os.Stdout.Write(content)
	}
}

func (r *Rotatee) Start() {
	log.WithFields(logrus.Fields{"Rotatee": r}).Debug("Start rotatee")
	pipeGroup := setupEventPipe(r.setting)
	pipeGroup.Start()
	pipeGroup.Broadcast(NewWriteTarget())
	teeLoop(pipeGroup)
}
