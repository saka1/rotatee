package main

type FormatSetter struct {
	format Format
}

func NewFormatSetter(format Format) FormatSetter {
	return FormatSetter{format}
}

func (fs FormatSetter) Run(in chan Event, out chan Event) {
	for {
		event, ok := <-in
		if !ok {
			close(out)
			return
		}
		switch event.eventType {
		case EventTypeInit, EventTypeChangeWriteTarget:
			event.format = fs.format
			out <- event
		default:
			out <- event
		}
	}
}
