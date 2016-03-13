package config

import (
	. "github.com/denbeigh2000/goi3bar"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type BuildFn func(options interface{}) (Producer, error)

type Config struct {
	Package string      `json:"package"`
	Name    string      `json:"name"`
	Options interface{} `json:"options"`
}

type ConfigSet struct {
	Entries  []Config `json:"entries"`
	Interval string   `json:"interval"`

	// Protects builders
	sync.Mutex
	builders map[string]BuildFn
}

func ReadConfigSet(path string) (cs ConfigSet, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	cs = ConfigSet{
		builders: make(map[string]BuildFn),
	}
	err = json.Unmarshal(data, &cs)
	if err != nil {
		return
	}

	return cs, nil
}

func (c *ConfigSet) Register(key string, builder BuildFn) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.builders[key]; ok {
		panic(fmt.Sprintf("Builder %s already exists, cannot reuse keys", key))
	}

	c.builders[key] = builder
}

func (c ConfigSet) Build() (bar *I3bar, err error) {
	keys := make([]string, len(c.Entries))
	interval, err := time.ParseDuration(c.Interval)
	if err != nil {
		return
	}

	bar = NewI3bar(interval)

	var producer Producer
	for i, e := range c.Entries {
		k := e.Package
		builder, ok := c.builders[k]
		if !ok {
			err = fmt.Errorf("Could not instantiate builder %v, unknown", k)
			return
		}

		producer, err = builder(e.Options)
		if err != nil {
			return
		}

		keys[i] = k
		bar.Register(e.Name, producer)
	}

	bar.Order(keys)

	return
}
