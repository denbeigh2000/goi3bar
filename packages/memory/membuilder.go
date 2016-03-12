package memory

import (
	. "github.com/denbeigh2000/goi3bar"

	"fmt"
	"time"
)

func MemoryBuilder(options interface{}) (Producer, error) {
	mc, ok := options.(memoryConfig)
	if !ok {
		return nil, fmt.Errorf("Could not format given data into memory config")
	}

	return &BaseProducer{
		Generator: Memory{},
		Name:      mc.Name,
		Interval:  mc.Refresh,
	}, nil
}

type memoryConfig struct {
	Name    string        `yaml:"name"`
	Refresh time.Duration `yaml:"refresh"`
}
