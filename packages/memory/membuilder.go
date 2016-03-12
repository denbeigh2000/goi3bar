package memory

import (
	. "github.com/denbeigh2000/goi3bar"

	"fmt"
	"time"
)

func Build(options interface{}) (Producer, error) {
	memMap, ok := options.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Could not format given data into memory config")
	}

	if _, ok := memMap["refresh"]; !ok {
		return nil, fmt.Errorf("Missing refresh field")
	}

	mb := memoryBuilder{
		Name:    memMap["name"].(string),
		Refresh: memMap["refresh"].(string),
	}

	return mb.Build()
}

func (m memoryBuilder) Build() (Producer, error) {
	interval, err := time.ParseDuration(m.Refresh)
	if err != nil {
		return nil, err
	}

	return &BaseProducer{
		Generator: Memory{},
		Name:      m.Name,
		Interval:  interval,
	}, nil
}

type memoryBuilder struct {
	Name    string `json:"name"`
	Refresh string `json:"refresh"`
}
