package network

import (
	"encoding/json"
	"fmt"
)

// guessJsonType makes some basic assertions about the given interface{} object
// (assumed to be directly decoded from JSON using json.Unmarshal) and uses it
// to narrow down the number of possible fields to one (or zero), to avoid
// calling our extremely expensive JSONReparse more than we need to.
func guessJsonType(i interface{}) string {
	iMap, ok := i.(map[string]interface{})

	fmt.Printf("i: %v\tiMap: %v\tok: %v\n", i, iMap, ok)
	if !ok {
		return ""
	}

	if _, ok := iMap["devices"]; ok {
		return "multi"
	}

	if _, ok := iMap["wireless"]; ok {
		return "wireless"
	}

	if _, ok := iMap["identifier"]; ok {
		return "basic"
	}

	return ""
}

// buildNetworkConfig makes some optimistic assertions about the type of the
// given interface, and tries to decode it into a know Devicer
func buildNetworkConfig(data []byte) (d Devicer, err error) {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	container := make(map[string]interface{})
	fmt.Println(string(data))
	err = json.Unmarshal(data, &container)
	if err != nil {
		return
	}

	var o interface{}
	switch guessJsonType(container) {
	case "multi":
		o = &multiDeviceConfig{}
	case "basic":
		o = &basicDeviceConfig{}
	case "wireless":
		o = &wirelessDeviceConfig{}
	case "":
		return nil, fmt.Errorf("Failed to decode interface{} into NetworkDevice")
	default:
		return nil, fmt.Errorf("Unknown NetworkDevice type: this should never happen")
	}

	err = json.Unmarshal(dataCopy, o)
	if err != nil {
		return
	}

	d = o.(Devicer)
	return
}
