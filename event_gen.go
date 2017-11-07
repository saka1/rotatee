package main

import (
	"github.com/sirupsen/logrus"
	"time"
)

type Timer struct {
	series Series
}

func NewTimer(series Series) Timer {
	return Timer{series}
}

func (t Timer) Run(in chan Event, out chan Event) {
	log.WithFields(logrus.Fields{"series": t.series}).Debug("Start rotate Timer")
	series := t.series
	series = series.Next()
	log.WithFields(logrus.Fields{"duration": series.Sub(time.Now())}).Debug("timer sleep")
	for {
		select {
		case event, ok := <-in:
			if !ok {
				close(out)
				return
			}
			out <- event
		case <-time.After(series.Sub(time.Now())):
			log.Debug("Rotation fired")
			event := NewWriteTarget()
			event.timestamp = series.Current()
			out <- event
			series = series.Next()
			log.WithFields(logrus.Fields{"duration": series.Sub(time.Now())}).Debug("Next sleep")
		}
	}
}
