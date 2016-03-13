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

	warnThresholdRaw, ok := memMap["warnThreshold"]
	if ok {
		warnThreshold, ok := warnThresholdRaw.(float64)
		if !ok {
			return nil, fmt.Errorf("Failed to parse given warnThreshold")
		}

		if warnThreshold > 100 || warnThreshold < 1 {
			return nil, fmt.Errorf("Invalid value for warnThreshold: %.0f", warnThreshold)
		}

		mb.WarnThreshold = int(warnThreshold)
	}

	critThresholdRaw, ok := memMap["critThreshold"]
	if ok {
		critThreshold, ok := critThresholdRaw.(float64)
		if !ok {
			return nil, fmt.Errorf("Failed to parse given critThreshold")
		}

		if critThreshold > 100 || critThreshold < 1 {
			return nil, fmt.Errorf("Invalid value for critThreshold: %.0f", critThreshold)
		}

		mb.CritThreshold = int(critThreshold)
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
	Name          string
	Refresh       string
	WarnThreshold int
	CritThreshold int
}
