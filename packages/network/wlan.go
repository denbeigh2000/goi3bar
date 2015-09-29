package network

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

const (
	essidRxStr      = "ESSID:\"(.*)\""
	strengthRxStr   = "Link Quality=(\\d+)/(\\d+)"
	notConnectedTpl = "%v not connected"
)

var (
	essidRx    = regexp.MustCompile(essidRxStr)
	strengthRx = regexp.MustCompile(strengthRxStr)
)

type WLANDevice struct {
	BasicNetworkDevice

	WarnThreshold int
	CritThreshold int

	// Signal strength as a number between 0 and 100
	strength int
	essid    string
}

func (d *WLANDevice) updateESSID(input string) error {
	matches := essidRx.FindStringSubmatch(input)
	if len(matches) < 2 {
		return fmt.Errorf("Couldn't match ESSID")
	}

	d.essid = matches[1]

	return nil
}

func (d *WLANDevice) updateStrength(input string) error {
	matches := strengthRx.FindStringSubmatch(input)
	if len(matches) < 3 {
		return fmt.Errorf("Couldn't match strength")
	}

	strengthNum, errN := strconv.Atoi(matches[1])
	strengthDenom, errD := strconv.Atoi(matches[1])
	if errN != nil {
		return errN
	}
	if errD != nil {
		return errD
	}

	d.strength = (strengthNum * 100) / strengthDenom

	return nil
}

func (d *WLANDevice) fetch() (string, error) {
	output, err := exec.Command("iwconfig", d.Identifier).Output()
	if err != nil {
		return "", err
	}
	outputS := string(output)

	return outputS, nil
}

func (d *WLANDevice) update() error {
	d.BasicNetworkDevice.Update()

	iwOut, err := d.fetch()
	if err != nil {
		return err
	}

	err = d.updateStrength(iwOut)
	if err != nil {
		return err
	}

	err = d.updateESSID(iwOut)
	if err != nil {
		return err
	}

	return nil
}

// Generate implements Generator
func (d *WLANDevice) Generate() ([]i3.Output, error) {
	err := d.update()
	if err != nil {
		return nil, err
	}

	if !d.connected {
		return []i3.Output{i3.Output{
			FullText: fmt.Sprintf(notConnectedTpl, d.Name),
			Color:    "#FF0000",
		}}, nil
	}

	var ip string
	if d.ip == nil {
		ip = "Acquiring IP"
	} else {
		ip = d.ip.String()
	}

	txt := fmt.Sprintf("%v: %v %v%% (%v)", d.Name, d.essid,
		d.strength, ip)

	var color string
	switch {
	case d.strength < d.CritThreshold:
		color = "#FF0000"
	case d.strength < d.WarnThreshold:
		color = "#FFA500"
	default:
		color = "#00FF00"
	}

	out := i3.Output{
		FullText:  txt,
		Color:     color,
		Separator: true,
	}

	return []i3.Output{out}, nil
}
