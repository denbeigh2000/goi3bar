package goi3bar

import (
	"fmt"
	"os"
	"time"
)

type Updater interface {
	Update(*Output) error
}

type Item struct {
	Updater

	Name    string
	Current *Output
	refresh time.Duration
	kill    chan struct{}
}

func (i Item) Start() {
	go i.loop()
}

func (i Item) Kill() {
	close(i.kill)
}

func (i *Item) loop() {
	t := time.NewTicker(i.refresh)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			err := i.Update(i.Current)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating output for %v: %v\n", i.Name, err)
				continue
			}

		case <-i.kill:
			return
		}
	}
}

func NewItem(name string, interval time.Duration, u Updater) *Item {
	item := Item{
		Updater: u,

		Name:    name,
		Current: &Output{Color: "#FFFFFF"},
		refresh: interval,
		kill:    make(chan struct{}),
	}

	err := u.Update(item.Current)
	if err != nil {
		panic(err)
	}

	return &item
}
