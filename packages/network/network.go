package network

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	notConnected = "Not connected"
	ethFormat    = "Connected: %v (%vMb/s)"
	ipAddrRxStr  = "inet\\s+((\\d{1,3}\\.){3}\\d{1,3})"
)

var (
	ipAddrRx = regexp.MustCompile(ipAddrRxStr)
)

// A BaseNetworkDevice describes a network device to be displayed on an i3bar
// TODO: Usage speed?
type NetworkDevice struct {
	// A friendly name to be used on the bar
	Name string

	// Name of the network interface for the corresponding device
	Identifier string

	ip net.IP
	// Link speed in kbits/sec
	speed     uint64
	connected bool
}

func (d *NetworkDevice) Update() error {
	output, err := exec.Command("ip", "addr", "show", d.Identifier).Output()
	if err != nil {
		return err
	}

	matches := ipAddrRx.FindStringSubmatch(string(output))
	d.ip = net.ParseIP(matches[1])

	// TODO: Bring crushing reality upon our users of their network speed
	d.speed = 1000000

	return nil
}

func (d *NetworkDevice) Generate() ([]i3.Output, error) {
	d.Update()

	speed := strconv.FormatUint(d.speed/1000, 10)

	text := fmt.Sprintf(ethFormat, d.ip.String(), speed)

	return []i3.Output{i3.Output{
		FullText: text,
		Color:    "#00FF00",
	}}, nil
}

// Network is a network applet, consisting of zero or more NetworkDevices,
// to be displayed on an i3bar
type NetworkDeviceCollection struct {
	// Map of interface name to relevant device
	Devices map[string]*NetworkDevice

	// Map of interface name to friendly name
	FriendlyNames map[string]string
}

// Generate implements Generator
func (c *NetworkDeviceCollection) Generate() ([]i3.Output, error) {
	return nil, nil
}
