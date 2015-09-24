package goi3bar

import (
	"fmt"
	"os"
	"time"
)

// A Generator generates content to put on an i3bar. Other functions will call
// Generate to create output for the i3bar.
type Generator interface {
	Generate() ([]Output, error)
}

// A Producer pushes content updates to the i3bar. It is responsible for
// generating its' own output in a timely manner
type Producer interface {
	Produce(kill <-chan struct{}) <-chan Update
}

// A BaseProducer is a simple Producer, which generates output at regular
// intervals using a Generator.
type BaseProducer struct {
	Generator

	Interval time.Duration
	Name     string
}

// A StaticGenerator is a simple Generator that returns the same Output each time.
type StaticGenerator []Output

// Generate implements Generator
func (g StaticGenerator) Generate() ([]Output, error) {
	return []Output(g), nil
}

// sendOutput is a helper function that waits up to p.Interval to send the
// given data down the given output channel, and abandons if it cannot.
func (p BaseProducer) sendOutput(out chan<- Update, data []Output) {
	select {
	case out <- Update{
		Key: p.Name,
		Out: data,
	}:
	case <-time.After(p.Interval):
	}
}

// Produce implements Producer. It creates a new value from the Generator every
// interval, and sends it down the provided channel
func (p *BaseProducer) Produce(kill <-chan struct{}) <-chan Update {
	out := make(chan Update)
	go func() {
		defer close(out)
		t := time.NewTicker(p.Interval)
		defer t.Stop()

		// Generate first pack to deliver information fast
		data, err := p.Generate()
		if err != nil {
			return
		}
		p.sendOutput(out, data)

		for {
			select {
			case <-t.C:
				data, err := p.Generate()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", p.Name, err)
					continue
				}

				// Attempt to send the output
				p.sendOutput(out, data)
			}
		}
	}()
	return out
}
