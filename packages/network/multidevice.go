package network

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"sync"
)

// because fuck waiting around for multiple networks to parse files
func (m *MultiDevice) concurrentUpdate() error {
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	for _, d := range m.Devices {
		wg.Add(1)
		go func(d NetworkDevice) {
			err := d.Update()
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(d)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for e := range errCh {
		return e
	}

	return nil
}

// MultiDevice is a Generator that manages a set of networks, and shows the
// most preferred network that is conneced
type MultiDevice struct {
	Devices map[string]NetworkDevice

	// Oredered list of keys describing preferred networks to show
	Preference []string
}

func (m MultiDevice) Generate() ([]i3.Output, error) {
	err := m.concurrentUpdate()
	if err != nil {
		return nil, err
	}

	for _, key := range m.Preference {
		if _, ok := m.Devices[key]; !ok {
			return nil, fmt.Errorf("Device %v doesn't exist", key)
		}

		device := m.Devices[key]
		if device.Connected() {
			return device.Generate()
		}
	}

	return []i3.Output{{
		FullText:  "No devices connected",
		Color:     i3.DefaultColors.Crit,
		Separator: true,
	}}, nil
}
