package goi3bar

import (
	"fmt"
	"os"
	"time"

	"github.com/denbeigh2000/goi3bar/util"
)

// A Generator generates content to put on an i3bar. Other functions will call
// Generate to create output for the i3bar.
// A Generator should define how, not when the output is built.
type Generator interface {
	Generate() ([]Output, error)
}

// A Producer pushes content updates to the i3bar. It is responsible for
// managing how often an item delivers its updates to the i3bar. These updates
// are usually generated using a Generator.
type Producer interface {
	Produce(kill <-chan struct{}) <-chan []Output
}

// A Clicker receives click events from the i3bar. If a registered Producer
// also implements Clicker, then that its' Click method will be called with
// the click event received from i3bar.
type Clicker interface {
	Click(ClickEvent) error
}

// A GeneratorClicker is both a Generator and a Clicker.
// It exists so we can have a concrete implementation built on BaseProducer
// that can also be used for plugins that implement Clicker
type GeneratorClicker interface {
	Generator
	Clicker
}

// A ProducerClicker is both a Producer and a Clicker.
// It exists so we can have a concrete implementation built on BaseProducer
// that can also be used for plugins that implement Clicker
type ProducerClicker interface {
	Producer
	Clicker
}

// A BaseProducer is a simple Producer, which generates output at regular
// intervals using a Generator.
type BaseProducer struct {
	Generator

	Interval time.Duration
	Name     string
}

// A BaseProducerClicker is like a BaseProducer, but it uses a GeneratorClicker
// instead of a Generator so it can also implement Clicker
type BaseProducerClicker struct {
	GeneratorClicker

	Interval time.Duration
	Name     string
}

func (p *BaseProducerClicker) Produce(kill <-chan struct{}) <-chan []Output {
	producer := &BaseProducer{
		Generator: p,

		Interval: p.Interval,
		Name:     p.Name,
	}

	return producer.Produce(kill)
}

// A StaticGenerator is a simple Generator that returns the same Output each time.
type StaticGenerator []Output

// Generate implements Generator
func (g StaticGenerator) Generate() ([]Output, error) {
	return []Output(g), nil
}

// sendOutput is a helper function that waits up to p.Interval to send the
// given data down the given output channel, and abandons if it cannot.
func (p BaseProducer) sendOutput(out chan<- []Output, data []Output, kill <-chan struct{}) {
	select {
	case out <- data:
	case <-time.After(p.Interval):
	case <-kill:
	}
}

// Produce implements Producer. It creates a new value from the Generator every
// interval, and sends it down the provided channel
func (p *BaseProducer) Produce(kill <-chan struct{}) <-chan []Output {
	out := make(chan []Output)
	go func() {
		defer close(out)
		t := util.NewTicker(p.Interval, true)
		defer t.Kill()

		for {
			select {
			case <-t.C:
				data, err := p.Generate()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", p.Name, err)
					data = []Output{{
						FullText:  fmt.Sprintf("ERROR from %v: %v", p.Name, err),
						Color:     DefaultColors.Crit,
						Separator: true,
					}}
				}

				// Attempt to send the output
				p.sendOutput(out, data, kill)

			case <-kill:
				return
			}
		}
	}()
	return out
}
