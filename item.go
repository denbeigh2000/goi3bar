package goi3bar

import (
	"fmt"
	"os"
	"time"
)

type Item struct {
	generate func(*Output) error

	Name    string
	current *Output
	refresh time.Duration
	kill    chan struct{}
}

func (i Item) Start() {
	go i.loop()
}

func (i Item) Kill() {
	close(i.kill)
}

func (i *Item) Get() *Output {
	//fmt.Printf("Get-ting %v\n", i.current)
	return i.current
}

func (i *Item) loop() {
	t := time.NewTicker(i.refresh)
	defer t.Stop()

	for {
		//fmt.Printf("Looping\n")
		select {
		case <-t.C:
			err := i.generate(i.current)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", i.Name, err)
				continue
			}

			//fmt.Printf("New value: %v\n", i.current)

		case <-i.kill:
			return
		}
	}
}

func NewItem(name string, interval time.Duration, fn func(*Output) error) *Item {
	item := Item{
		generate: fn,

		Name:    name,
		current: &Output{Color: "#FFFFFF"},
		refresh: interval,
		kill:    make(chan struct{}),
	}

	err := fn(item.current)
	if err != nil {
		panic(err)
	}

	return &item
}
