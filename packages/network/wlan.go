package network

import (
	i3 "github.com/denbeigh2000/goi3bar"
)

type WLANDevice struct {
	NetworkDevice

	strength int
	essid    string
}

func (d *WLANDevice) update() error {
	return nil
}

// Generate implements Generator
func (d *WLANDevice) Generate() ([]i3.Output, error) {
	err := d.update()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
