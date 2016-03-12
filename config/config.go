package config

import (
	. "github.com/denbeigh2000/goi3bar"

	"gopkg.in/yaml.v2"

	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Builder interface {
	Build(options interface{}) (Producer, error)
}

type Config struct {
	Package string      `yaml:"package"`
	Options interface{} `yaml:"options"`
}

type ConfigSet struct {
	Entries  []Config `yaml:"entries"`
	Interval time.Duration

	// Protects builders
	sync.Mutex
	builders map[string]Builder
}

func ReadConfigFile(path string) (I3bar, error) {
	return I3bar{}, nil
}

func readConfigSet(path string) (cs ConfigSet, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, cs)
	if err != nil {
		return
	}

	return cs, nil
}

func (c *ConfigSet) Register(key string, builder Builder) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.builders[key]; ok {
		panic(fmt.Sprintf("Builder %s already exists, cannot reuse keys", key))
	}

	c.builders[key] = builder
}

func (c ConfigSet) Build() (bar I3bar, err error) {
	keys := make([]string, len(c.Entries))

	var producer Producer
	for i, e := range c.Entries {
		k := e.Package
		builder := c.builders[k]

		producer, err = builder.Build(e.Options)
		if err != nil {
			return
		}

		keys[i] = k
		bar.Register(k, producer)
	}

	bar.Order(keys)

	return
}
