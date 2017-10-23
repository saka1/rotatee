package main

type FormatEvaluator struct {
	format string
}

func newFormatEval(format string) FormatEvaluator {
	return FormatEvaluator{format}
}

func (fe *FormatEvaluator) Start(in chan Event, out chan Event) {
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
			event.fileName = fe.format //TODO eval
			out <- event
		default:
			out <- event
		}
	}
}
