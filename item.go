package goi3bar

import (
	"fmt"
	"os"
	"time"
)

type Generator interface {
	Generate() ([]Output, error)
}

type Item struct {
	Generator

	Name    string
	refresh time.Duration
}

func (i *Item) sendItems(ch chan<- Update, items []Output) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		ch <- Update{
			Key: i.Name,
			Out: items,
		}

		close(done)
	}()

	return done
}

func (i *Item) sendOutput(out chan<- Update) error {
	outputs, err := i.Generate()
	if err != nil {
		return err
	}

	// Try to asynchronously send the output, if it's time for another output pack, abandon it
	go func() {
		select {
		case <-i.sendItems(out, outputs):
			return
		case <-time.After(i.refresh):
			return
		}
	}()

	return nil
}

func (i *Item) Produce(out chan<- Update, kill <-chan struct{}) {
	t := time.NewTicker(i.refresh)
	defer t.Stop()

	err := i.sendOutput(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", i.Name, err)
	}

	for {
		select {
		case <-t.C:
			err := i.sendOutput(out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", i.Name, err)
				continue
			}

		}
	}
}

type BaseProducer struct {
	Generator

	Interval time.Duration
	Name     string
}

func (p BaseProducer) sendItems(ch chan<- Update, items []Output) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		ch <- Update{
			Key: p.Name,
			Out: items,
		}

		close(done)
	}()

	return done
}

func (p BaseProducer) sendOutput(out chan<- Update) error {
	outputs, err := p.Generate()
	if err != nil {
		return err
	}

	// Try to asynchronously send the output, if it's time for another output pack, abandon it
	go func() {
		select {
		case <-p.sendItems(out, outputs):
			return
		case <-time.After(p.Interval):
			return
		}
	}()

	return nil
}

func (p *BaseProducer) Produce(out chan<- Update, kill <-chan struct{}) {
	t := time.NewTicker(p.Interval)
	defer t.Stop()

	err := p.sendOutput(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", p.Name, err)
	}

	for {
		select {
		case <-t.C:
			err := p.sendOutput(out)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", p.Name, err)
				continue
			}

		}
	}
}

func NewItem(interval time.Duration, g Generator) *Item {
	item := Item{
		Generator: g,

		refresh: interval,
	}

	return &item
}
