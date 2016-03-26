package config

import (
	. "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/util"

	"fmt"
	"sync"
	"time"
)

var (
	lock     sync.RWMutex
	builders = make(map[string]Builder, 0)
)

// Builder is the interface that must be implemented by plugins that wish to be
// configurable with I3bar's JSON configuration. Its' Build() method will be
// called exactly once on start with the given config. The Builder is strongly
// advised (though not required) to take advantage of the Config's ParseConfig()
// method, which parses the JSON options struct into a struct of their choosing.
// It behaves exactly like json.Unmarshal.
type Builder interface {
	Build(Config) (Producer, error)
}

// Config for one individual plugin instance
type Config struct {
	Package string      `json:"package"`
	Name    string      `json:"name"`
	Options interface{} `json:"options"`
}

// ParseConfig is the point where your plugin's JSON config subtree will be
// parsed. Call this function with a pointer to your JSON-annotated config
// struct type in here, and it will behave as you expect it to.
func (c Config) ParseConfig(i interface{}) error {
	return util.JSONReparse(c.Options, i)
}

// ConfigSet represents an entire JSON config file. If all goes well, it
// creates an I3bar.
type ConfigSet struct {
	Entries  []Config `json:"entries"`
	Interval string   `json:"interval"`
}

// Build() constructs an I3bar from its internal configuration. The returned
// I3bar will not have had its' Start() method called.
func (c ConfigSet) Build() (bar *I3bar, err error) {
	keys := make([]string, len(c.Entries))
	interval, err := time.ParseDuration(c.Interval)
	if err != nil {
		return
	}

	bar = NewI3bar(interval)

	lock.RLock()
	defer lock.RUnlock()

	var producer Producer
	for i, e := range c.Entries {
		k := e.Package
		builder, ok := builders[k]
		if !ok {
			err = fmt.Errorf("Could not instantiate builder %v, unknown", k)
			return
		}

		producer, err = builder.Build(e)
		if err != nil {
			return
		}

		keys[i] = k
		bar.Register(e.Name, producer)
	}

	bar.Order(keys)

	return
}

// Register is the registration point of your plugin to the Configuration
// namespace. Use a unique key, and your Builder will be called for entries
// with that package key
func Register(key string, builder Builder) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := builders[key]
	if ok {
		panic(fmt.Sprintf("Builder %v already exists", key))
	}

	builders[key] = builder
}
