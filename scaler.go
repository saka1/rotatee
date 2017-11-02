package main

type Scaler struct {
	limit int64
}

func NewScaler(limit int64) Scaler {
	return Scaler{limit}
}

func (s Scaler) Run(in chan Event, out chan Event) {
	var count int64 = 0
	for {
		event, ok := <-in
		if !ok {
			close(out)
			return
		}
		switch event.eventType {
		case EVENT_TYPE_PAYLOAD:
			unprocessedPayload := event.payload
			for count+int64(len(unprocessedPayload)) > s.limit {
				acceptable := s.limit - count
				currentPayload := make([]byte, acceptable)
				copy(currentPayload, unprocessedPayload)
				out <- NewPayload(currentPayload)
				// interleave rolling event
				out <- NewWriteTarget()
				unprocessedPayload = unprocessedPayload[acceptable:]
				count = 0
			}
			lastPayload := make([]byte, len(unprocessedPayload))
			copy(lastPayload, unprocessedPayload) //TODO optimize copy?
			out <- NewPayload(lastPayload)
			count = int64(len(lastPayload))
		default:
			out <- event
		}
	}
}
