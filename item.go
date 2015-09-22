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
	out     chan<- Update
	refresh time.Duration
	kill    chan struct{}
}

func (i Item) Start() {
	if i.out == nil || i.kill == nil {
		panic("Item must be registered before starting")
	}

	go i.loop()
}

func (i *Item) sendItems(items []Output) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		for _, item := range items {
			i.out <- Update{
				Key: i.Name,
				Out: item,
			}
		}

		close(done)
	}()

	return done
}

func (i *Item) sendOutput() error {
	outputs, err := i.Generate()
	if err != nil {
		return err
	}

	// Try to asynchronously send the output, if it's time for another output pack, abandon it
	go func() {
		select {
		case <-i.sendItems(outputs):
			return
		case <-time.After(i.refresh):
			return
		}
	}()

	return nil
}

func (i *Item) loop() {
	t := time.NewTicker(i.refresh)
	defer t.Stop()

	err := i.sendOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", i.Name, err)
	}

	for {
		select {
		case <-t.C:
			err := i.sendOutput()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", i.Name, err)
				continue
			}

		case <-i.kill:
			return
		}
	}
}

func NewItem(name string, interval time.Duration, g Generator) *Item {
	item := Item{
		Generator: g,

		Name:    name,
		out:     nil,
		refresh: interval,
		kill:    nil,
	}

	return &item
}
