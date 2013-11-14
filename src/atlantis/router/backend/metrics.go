package backend

type ServerMetrics struct {
	RequestsInFlight uint32
	RequestsServiced uint64
}

func NewServerMetrics() ServerMetrics {
	return ServerMetrics{
		RequestsInFlight: 0,
		RequestsServiced: 0,
	}
}

func (s *ServerMetrics) RequestStart() {
	s.RequestsInFlight++
	s.RequestsServiced++
}

func (s *ServerMetrics) RequestDone() {
	s.RequestsInFlight--
}

func (s *ServerMetrics) Cost() uint32 {
	return s.RequestsInFlight
}
