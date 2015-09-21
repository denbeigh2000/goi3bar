package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	timeFormat "github.com/jehiah/go-strftime"
)

const (
	intro        = "{ \"version\": 1 }\n"
	formatString = "%a %d-%b-%y %I:%M:%S"
)

type Color struct {
	r uint32
	g uint32
	b uint32
}

type Output struct {
	Align     string `json:"align,omitEmpty"`
	Color     string `json:"color,omitEmpty"`
	FullText  string `json:"full_text"`
	Instance  string `json:"instance,omitEmpty"`
	MinWidth  string `json:"min_width,omitEmpty"`
	Name      string `json:"name,omitEmpty"`
	ShortText string `json:"short_text,omitEmpty"`
	Separator bool   `json:"separator"`
	Urgent    bool   `json:"urgent"`
}

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

func output(ch <-chan []*Output) {
	fmt.Fprintf(os.Stdout, intro)
	fmt.Fprintf(os.Stdout, "[\n")

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	for o := range ch {
		_, err := buf.Write([]byte("["))
		if err != nil {
			panic(err)
		}

		isFirst := true
		for _, item := range o {
			if !isFirst {
				_, err := buf.Write([]byte(","))
				if err != nil {
					panic(err)
				}
			} else {
				isFirst = false
			}

			//fmt.Println(item.FullText)
			err := enc.Encode(item)
			if err != nil {
				panic(err)
			}
		}

		_, err = buf.Write([]byte("],\n"))
		if err != nil {
			panic(err)
		}

		io.Copy(os.Stdout, &buf)
	}
}

type I3bar struct {
	items    map[string]*Item
	order    []string
	interval time.Duration
	json     chan []*Output
	kill     chan struct{}
}

func NewI3bar(update time.Duration) *I3bar {
	return &I3bar{
		items:    make(map[string]*Item),
		order:    make([]string, 0),
		interval: update,
		json:     make(chan []*Output),
		kill:     make(chan struct{}),
	}
}

func (i I3bar) Start() {
	go i.loop()
}

func (i I3bar) Kill() {
	close(i.kill)
}

func (i *I3bar) Register(key string, item *Item) {
	_, ok := i.items[key]
	if ok {
		panic(fmt.Sprintf("Key %v exists", key))
	}

	i.items[key] = item
	i.order = append(i.order, key)

	item.Start()
}

func (i *I3bar) Order(keys []string) error {
	if len(keys) != len(i.items) {
		return fmt.Errorf("Number of keys must equal number of items, expected %v got %v",
			len(i.items), len(keys))
	}

	for _, k := range keys {
		if _, ok := i.items[k]; !ok {
			return fmt.Errorf("Key not present: %v", k)
		}
	}

	i.order = keys
	return nil
}

func (i *I3bar) collect() []*Output {
	items := make([]*Output, len(i.items))

	for idx, k := range i.order {
		item, ok := i.items[k]
		if !ok {
			panic(fmt.Sprintf("Missing key %v", k))
		}

		items[idx] = item.Get()
	}
	return items
}

func (i *I3bar) loop() {
	defer close(i.json)

	t := time.NewTicker(i.interval)
	defer t.Stop()
	defer func() {
		for _, item := range i.items {
			item.Kill()
		}
	}()

	go output(i.json)

	for {
		select {
		case <-t.C:
			items := i.collect()

			select {
			case <-i.kill:
				return
			case i.json <- items:
				//fmt.Println(items[0].FullText)
				continue
			}
		case <-i.kill:
			return
		}
	}
}

func formatTime() string {
	time := time.Now()

	return timeFormat.Format(formatString, time)
}

func timeOutput(o *Output) error {
	o.FullText = formatTime()
	return nil
}

func main() {
	clock := NewItem("time", 1*time.Second, timeOutput)

	bar := NewI3bar(1 * time.Second)

	bar.Register("time", clock)

	bar.Start()
	defer bar.Kill()

	<-time.After(20 * time.Second)
}
