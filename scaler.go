package main

type Scaler struct {
	limit uint64
}

func newScaler(limit uint64) Scaler {
	return Scaler{limit}
}

func (s Scaler) Run(in chan Event, out chan Event) {
	var count uint64 = 0
	for {
		event, ok := <-in
		if !ok {
			close(out)
			return
		}
		switch event.eventType {
		case EVENT_TYPE_PAYLOAD:
			if count+uint64(len(event.payload)) > s.limit {
				// send head
				acceptable := s.limit - count
				firstPayload := make([]byte, acceptable)
				copy(firstPayload, event.payload)
				out <- NewPayload(firstPayload)
				// interleave rolling event
				out <- NewWriteTarget()
				// send rest
				secondPayload := make([]byte, uint64(len(event.payload))-acceptable)
				copy(secondPayload, event.payload[acceptable:])
				out <- NewPayload(secondPayload)
				// update counter
				count = uint64(len(secondPayload))
			} else {
				out <- event
				count += uint64(len(event.payload))
			}
		default:
			out <- event
		}
	}
}
