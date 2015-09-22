package goi3bar

func NewMultiProducer(m map[string]Producer) MultiProducer {
	return MultiProducer{m}
}

type MultiProducer struct {
	producers map[string]Producer
}

func (m MultiProducer) Produce(out chan<- Update, kill <-chan struct{}) {
	for _, p := range m.producers {
		go p.Produce(out, kill)
	}
}

func (m MultiProducer) MultiRegister(r Registerer) {
	for k, p := range m.producers {
		r.Register(k, p)
	}
}
