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
	// set default value
	if setting.stdin == nil {
		setting.stdin = os.Stdin
	}
	if setting.stdout == nil {
		setting.stdout = os.Stdout
	}
	return &Rotatee{setting: setting}
}

type RotateeSetting struct {
	stdin       io.Reader
	stdout      io.Writer
	args        []string
	verbose     bool
	maxFileSize int64
	historySize int
	appendMode  bool
}

func setupEventPipe(setting RotateeSetting) EventPipeGroup {
	pipeGroup := NewEventPipeGroup()
	for _, arg := range setting.args {
		pipe := NewEventPipe()
		if setting.maxFileSize != 0 {
			pipe.Add(NewScaler(setting.maxFileSize))
		}
		format := Format(arg)
		pipe.Add(NewTimer(DetectSeries(arg, time.Now())))
		pipe.Add(NewFormatSetter(format))
		pipe.Add(NewRoller(format, setting))
		pipeGroup.Add(&pipe)
	}
	return pipeGroup
}

func (r *Rotatee) teeLoop(pipeGroup *EventPipeGroup) {
	reader := r.setting.stdin
	readBuf := make([]byte, 1024)
	for {
		length, err := reader.Read(readBuf)
		if err != nil {
			if err == io.EOF {
				pipeGroup.Stop()
				return
			}
			log.Panic("Writer goroutine IO failed")
		}
		// do copy because 'content' is shared among goroutine(s)
		content := make([]byte, length)
		copy(content, readBuf[:length])
		pipeGroup.Broadcast(NewPayload(content))
		// Also, write to stdout
		r.setting.stdout.Write(content)
	}
}

func (r *Rotatee) Start() {
	log.WithFields(logrus.Fields{"Rotatee": r}).Debug("Start rotatee")
	pipeGroup := setupEventPipe(r.setting)
	pipeGroup.Start()
	pipeGroup.Broadcast(NewInit())
	r.teeLoop(&pipeGroup)
}
