package goi3bar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	intro        = "{ \"version\": 1 }\n"
	formatString = "%a %d-%b-%y %I:%M:%S"
)

type Producer interface {
	Produce(out chan<- Update, kill <-chan struct{})
}

type Registerer interface {
	Register(key string, p Producer)
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

func output(ch <-chan []Output) {
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

type Update struct {
	Key string
	Out []Output
}

type I3bar struct {
	producers map[string]Producer
	values    map[string][]Output
	order     []string
	interval  time.Duration
	in        chan Update
	json      chan []Output
	kill      chan struct{}
}

func NewI3bar(update time.Duration) *I3bar {
	return &I3bar{
		producers: make(map[string]Producer),
		order:     make([]string, 0),
		interval:  update,
		in:        make(chan Update),
		json:      make(chan []Output),
		kill:      make(chan struct{}),
		values:    make(map[string][]Output),
	}
}

func (i *I3bar) Start() {
	go i.loop()
}

func (i I3bar) Kill() {
	close(i.kill)
}

func (i *I3bar) Register(key string, p Producer) {
	_, ok := i.producers[key]
	if ok {
		panic(fmt.Sprintf("Producer %v exists", key))
	}

	i.producers[key] = p
	i.values[key] = nil
	i.order = append(i.order, key)

	go p.Produce(i.in, i.kill)
}

func (i *I3bar) Order(keys []string) error {
	if len(keys) != len(i.producers) {
		return fmt.Errorf("Number of keys must equal number of items, expected %v got %v",
			len(i.producers), len(keys))
	}

	for _, k := range keys {
		if _, ok := i.producers[k]; !ok {
			return fmt.Errorf("Producer not present: %v", k)
		}
	}

	i.order = keys
	return nil
}

func (i *I3bar) collect() []Output {
	items := make([]Output, 0)

	for _, k := range i.order {
		v, ok := i.values[k]
		if !ok {
			panic(fmt.Sprintf("Missing key %v", k))
		}

		for _, out := range v {
			items = append(items, out)
		}
	}

	return items
}

func (i *I3bar) loop() {
	defer close(i.json)

	t := time.NewTicker(i.interval)
	defer t.Stop()

	go output(i.json)

	for {
		select {
		case update := <-i.in:
			i.values[update.Key] = update.Out
		case <-t.C:
			items := i.collect()

			select {
			case <-i.kill:
				return
			case i.json <- items:
				continue
			}
		case <-i.kill:
			return
		}
	}
}
