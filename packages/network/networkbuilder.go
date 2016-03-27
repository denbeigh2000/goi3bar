package network

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/config"

	"fmt"
	"time"
)

// generalConfig represents the outermost configuration layer for any network
// device. It only has one property: refresh interval, and then can contain
// arbitrary device configs
type generalConfig struct {
	Interval string      `json:"interval"`
	Config   interface{} `json:"config"`
}

// Devicer is an interface used within this package to represent a configuration struct
// that is able to produce a Generator (representing one or more network devices)
type Devicer interface {
	Device() (i3.Generator, error)
}

// basicDeviceConfig represents a JSON configuration for a simple network device
// (such as ethernet)
type basicDeviceConfig struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

// Device() implements Devicer
func (c basicDeviceConfig) Device() (i3.Generator, error) {
	return &BasicNetworkDevice{
		Name:       c.Name,
		Identifier: c.Identifier,
	}, nil
}

// wirelessDeviceConfig represents a JSON configuration for a WLAN device
type wirelessDeviceConfig struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Wireless   struct {
		WarnThreshold int `json:"warn_threshold"`
		CritThreshold int `json:"crit_threshold"`
	} `json:"wireless"`
}

// Device() implements Devicer
func (c wirelessDeviceConfig) Device() (i3.Generator, error) {
	return &WLANDevice{
		BasicNetworkDevice: BasicNetworkDevice{
			Name:       c.Name,
			Identifier: c.Identifier,
		},
		WarnThreshold: c.Wireless.WarnThreshold,
		CritThreshold: c.Wireless.CritThreshold,
	}, nil
}

// multiDeviceConfig represents a JSON configuration for a MultiDevice
type multiDeviceConfig struct {
	Devices    map[string]interface{} `json:"devices"`
	Preference []string               `json:"preference"`
}

// Device() implements Devicer
func (c multiDeviceConfig) Device() (i3.Generator, error) {
	if len(c.Devices) != len(c.Preference) {
		return MultiDevice{}, fmt.Errorf("Number of network devices and number of preferences must be equal")
	}

	devices := make(map[string]NetworkDevice)
	for k, v := range c.Devices {
		// Prevent erroneous and recursive definitions before expensive operations
		switch guessJsonType(v) {
		case "basic", "wireless":
		default:
			return MultiDevice{}, fmt.Errorf("Can't have recursive multi device configs")
		}

		// Attempt to determine which network device we are given
		d, err := buildNetworkConfig(v)
		if err != nil {
			return MultiDevice{}, err
		}

		if _, ok := devices[k]; ok {
			return MultiDevice{}, fmt.Errorf("Duplicate key %v", k)
		}

		// Produce a NetworkDevice and store it in our map
		generator, err := d.Device()
		if err != nil {
			return MultiDevice{}, err
		}

		switch t := generator.(type) {
		case *BasicNetworkDevice:
			devices[k] = t
		case *WLANDevice:
			devices[k] = t
		default:
			panic("Unsupported NetworkDevice type found while parsing config, second-time reflection. This should never happen")
		}
	}

	return MultiDevice{
		Devices:    devices,
		Preference: c.Preference,
	}, nil
}

// networkBuilder provides the entry point for config to produce a Network plugin
// from a config file entry
type networkBuilder struct{}

// Build implements config.Builder
func (b networkBuilder) Build(data config.Config) (p i3.Producer, err error) {
	var c generalConfig
	err = data.ParseConfig(&c)
	if err != nil {
		return
	}

	interval, err := time.ParseDuration(c.Interval)
	if err != nil {
		return
	}

	conf, err := buildNetworkConfig(c.Config)
	if err != nil {
		return
	}

	generator, err := conf.Device()
	if err != nil {
		return
	}

	p = &i3.BaseProducer{
		Generator: generator,
		Interval:  interval,
		Name:      "network",
	}

	return
}

func init() {
	config.Register("network", networkBuilder{})
}
