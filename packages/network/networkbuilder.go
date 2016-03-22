package network

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"time"
)

type basicDeviceConfig struct {
	name       string
	identifier string
}

func (c basicDeviceConfig) toDevice() *BasicNetworkDevice {
	return &BasicNetworkDevice{
		Name:       c.name,
		Identifier: c.identifier,
	}
}

func (c wirelessDeviceConfig) toDevice() *WLANDevice {
	return &WLANDevice{
		BasicNetworkDevice: BasicNetworkDevice{
			Name:       c.name,
			Identifier: c.identifier,
		},
		WarnThreshold: c.wireless.warnThreshold,
		CritThreshold: c.wireless.critThreshold,
	}
}

type wirelessDeviceConfig struct {
	name       string
	identifier string
	wireless   struct {
		warnThreshold int
		critThreshold int
	}
}

type multiDeviceConfig struct {
	devices    map[string]interface{}
	preference []string
}

func Build(options interface{}) (i3.Producer, error) {
	config, ok := options.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Couldn't decode given config %v", options)
	}

	intervalInt, ok := config["interval"]
	if !ok {
		return nil, fmt.Errorf("Missing interval field")
	}
	intervalStr, ok := intervalInt.(string)
	if !ok {
		return nil, fmt.Errorf("Couldn't decode interval to string")
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}

	networkOpts, ok := config["deviceConfig"]
	if !ok {
		return nil, fmt.Errorf("No device configuration section")
	}

	var g i3.Generator
	switch t := networkOpts.(type) {
	case multiDeviceConfig:
		g, err = buildMultiDevice(t)
		if err != nil {
			return nil, err
		}
	case wirelessDeviceConfig:
		g = t.toDevice()
	case basicDeviceConfig:
		g = t.toDevice()
	default:
		return nil, fmt.Errorf("Couldn't cast given type %v to a network device config", t)
	}

	return &i3.BaseProducer{
		Interval:  interval,
		Generator: g,
		Name:      "network",
	}, nil
}

func buildNetworkDevice(optRaw interface{}) (NetworkDevice, error) {
	options, ok := optRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid options struct given: %v", optRaw)
	}

}

func buildMultiDevice(options multiDeviceConfig) (i3.Generator, error) {
	if len(options.devices) != len(options.preference) {
		return nil, fmt.Errorf("Number of network devices and number of preferences must be equal")
	}

	devices := make(map[string]NetworkDevice)
	for k, v := range options.devices {
		var d NetworkDevice
		switch device := v.(type) {
		case basicDeviceConfig:
			d = device.toDevice()
		case wirelessDeviceConfig:
			d = device.toDevice()
		case multiDeviceConfig:
			return nil, fmt.Errorf("Can't have recursive multi device configs")
		default:
			return nil, fmt.Errorf("Unrecognised device type %v", device)
		}

		if _, ok := devices[k]; ok {
			return nil, fmt.Errorf("Duplicate key %v", k)
		}

		devices[k] = d
	}

	return MultiDevice{
		Devices:    devices,
		Preference: options.preference,
	}, nil
}
