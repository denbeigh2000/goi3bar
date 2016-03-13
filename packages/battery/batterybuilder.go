package battery

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"time"
)

func Build(optsRaw interface{}) (i3.Producer, error) {
	options, ok := optsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to parse battery JSON")
	}

	idRaw, ok := options["identifier"]
	if !ok {
		return nil, fmt.Errorf("Battery missing identifier")
	}
	id, ok := idRaw.(string)
	if !ok {
		return nil, fmt.Errorf("Identifier not a valid string")
	}

	intervalRaw, ok := options["interval"]
	if !ok {
		return nil, fmt.Errorf("Battery missing interval")
	}
	intervalStr, ok := intervalRaw.(string)
	if !ok {
		return nil, fmt.Errorf("Interval not a valid string")
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}

	bat := Battery{
		Identifier: id,
	}

	if nameRaw, ok := options["name"]; ok {
		if name, ok := nameRaw.(string); ok {
			bat.Name = name
		} else {
			return nil, fmt.Errorf("Failed to parse battery name")
		}
	}

	if warnThresholdRaw, ok := options["warnThreshold"]; ok {
		warnThreshold, ok := warnThresholdRaw.(float64)
		if !ok {
			return nil, fmt.Errorf("Failed to parse given warnThreshold")
		}

		if warnThreshold > 100 || warnThreshold < 1 {
			return nil, fmt.Errorf("Invalid value for warnThreshold: %.0f", warnThreshold)
		}

		bat.WarnThreshold = int(warnThreshold)
	}

	if critThresholdRaw, ok := options["critThreshold"]; ok {
		critThreshold, ok := critThresholdRaw.(float64)
		if !ok {
			return nil, fmt.Errorf("Failed to parse given critThreshold")
		}

		if critThreshold > 100 || critThreshold < 1 {
			return nil, fmt.Errorf("Invalid value for critThreshold: %.0f", critThreshold)
		}

		bat.CritThreshold = int(critThreshold)
	}

	return &i3.BaseProducer{
		Generator: &bat,
		Interval:  interval,
		Name:      "bat",
	}, nil
}
