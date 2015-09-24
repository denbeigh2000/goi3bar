package network

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"strconv"
)

const ethFormat = "Ethernet connected: %v (%v/s)"

type EthernetDevice struct {
	NetworkDevice

	// Link speed in kbits/sec
	speed uint64
}

func (d *EthernetDevice) Update() error {
	d.NetworkDevice.Update()

	// TODO: Bring crushing reality upon our users of their network speed
	d.speed = 1000000

	return nil
}

func (d *EthernetDevice) Generate() ([]i3.Output, error) {
	err := d.Update()
	if err != nil {
		return nil, err
	}

	speed := strconv.FormatUint(d.speed/1000, 10) + "Mb"

	text := fmt.Sprintf(ethFormat, d.ip.String(), speed)

	return []i3.Output{i3.Output{
		FullText: text,
		Color:    "#00FF00",
	}}, nil
}
