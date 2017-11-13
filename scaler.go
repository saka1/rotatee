package main

type Scaler struct {
	limit int64
	count int64
}

func NewScaler(limit int64) *Scaler {
	return &Scaler{limit, 0}
}

func (s *Scaler) Run(in chan Event, out chan Event) {
	for {
		event, ok := <-in
		if !ok {
			close(out)
			return
		}
		switch event.eventType {
		case EventTypePayload:
			unprocessedPayload := event.payload
			for s.count+int64(len(unprocessedPayload)) > s.limit {
				acceptable := s.limit - s.count
				currentPayload := make([]byte, acceptable)
				copy(currentPayload, unprocessedPayload)
				out <- NewPayload(currentPayload)
				// interleave rolling event
				out <- NewWriteTarget()
				unprocessedPayload = unprocessedPayload[acceptable:]
				s.count = 0
			}
			lastPayload := make([]byte, len(unprocessedPayload))
			copy(lastPayload, unprocessedPayload) //TODO optimize copy?
			out <- NewPayload(lastPayload)
			s.count += int64(len(lastPayload))
		default:
			out <- event
		}
	}
}
