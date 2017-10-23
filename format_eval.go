// FormatEvaluator capture specific message, add fileName from its timestamp.
package main

import (
	"github.com/jehiah/go-strftime"
)

type FormatEvaluator struct {
	format string
}

func NewFormatEval(format string) FormatEvaluator {
	return FormatEvaluator{format}
}

func (fe FormatEvaluator) Run(in chan Event, out chan Event) {
	for {
		event, ok := <-in
		if !ok {
			close(out)
			return
		}
		switch event.eventType {
		case EVENT_TYPE_PAYLOAD:
			out <- event // pass-through
		case EVENT_TYPE_CHANGE_WRITE_TARGET:
			event.fileName = strftime.Format(fe.format, event.timestamp)
			out <- event
		default:
			out <- event
		}
	}
}
