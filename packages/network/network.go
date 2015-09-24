package network

import (
	i3 "github.com/denbeigh2000/goi3bar"
)

// A NetworkDevice describes a network device to be displayed on an i3bar
type NetworkDevice struct {
	// A friendly name to be used on the bar
	Name string

	// Name of the network interface for the corresponding device
	Identifier string
}

// Network is a network applet, consisting of zero or more NetworkDevices,
// to be displayed on an i3bar
type Network struct {
	Devices map[string]*NetworkDevice
}

// Generate implements Generator
func (n *Network) Generate() ([]i3.Output, error) {
	return nil, nil
}
