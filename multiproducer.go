package goi3bar

// MultiProducer is a simple Producer that groups multiple Producers.
type MultiProducer struct {
	producers map[string]Producer
}

// NewMultiProducer creates a new MultiProducer
func NewMultiProducer(m map[string]Producer) MultiProducer {
	return MultiProducer{m}
}

// Produce implements Producer
func (m MultiProducer) Produce(out chan<- Update, kill <-chan struct{}) {
	for _, p := range m.producers {
		go p.Produce(out, kill)
	}
}

// MultiRegister takes a Registerer and uses it to register all of its'
// Producers
func (m MultiProducer) MultiRegister(r Registerer) {
	for k, p := range m.producers {
		r.Register(k, p)
	}
}
