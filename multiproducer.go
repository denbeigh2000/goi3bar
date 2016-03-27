package goi3bar

import "sync"

// MultiProducer is a simple Producer that groups multiple Producers.
type MultiProducer struct {
	producers map[string]Producer
}

// NewMultiProducer creates a new MultiProducer
func NewMultiProducer(m map[string]Producer) MultiProducer {
	return MultiProducer{m}
}

// Produce implements Producer
func (m MultiProducer) Produce(kill <-chan struct{}) <-chan []Output {
	out := make(chan []Output)
	wg := sync.WaitGroup{}
	for _, p := range m.producers {
		wg.Add(1)
		go func(p Producer) {
			defer wg.Done()
			ch := p.Produce(kill)
			for x := range ch {
				out <- x
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// MultiRegister takes a registerer and uses it to register all of its'
// Producers
func (m MultiProducer) MultiRegister(r registerer) {
	for k, p := range m.producers {
		r.Register(k, p)
	}
}
