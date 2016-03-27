package battery

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"io/ioutil"
	"sort"
	"strings"
)

var (
	// Lowercase map of power supply identifiers we don't want to include
	prohibitedItems = map[string]struct{}{
		"ac": {},
	}
)

// Creates a list of interfaces that can be uesd by BatteryDiscover
func findBatteries() ([]string, error) {
	dirs, err := ioutil.ReadDir(BaseBatteryPath)
	if err != nil {
		return nil, err
	}

	interfaces := make([]string, 0)
	for _, d := range dirs {
		name := d.Name()
		if _, ok := prohibitedItems[strings.ToLower(name)]; ok {
			// We are not interested in this power supply
			// (probably because it is a charger)
			continue
		}

		interfaces = append(interfaces, name)
	}

	return interfaces, nil
}

// Discovers batteries on your machine, and returns an OrderedMultiGenerator for each.
// Include a map of lowercase interface names to friendly names to display the
// batteries with friendly names.
func BatteryDiscover(names map[string]string, WarnThreshold, CritThreshold int) (i3.Generator, error) {
	interfaces, err := findBatteries()
	if err != nil {
		return nil, err
	}

	// Show a useful message if no batteries were found
	if len(interfaces) == 0 {
		return i3.StaticGenerator(
			[]i3.Output{{
				FullText:  "No battery found",
				Separator: true,
				Color:     "#FF0000",
			}},
		), nil
	}

	batteries := make(map[string]i3.Generator, len(interfaces))
	order := make([]string, len(interfaces))

	for i, in := range interfaces {
		var name string
		if n, ok := names[in]; ok {
			name = n
		} else {
			name = in
		}

		batteries[name] = &Battery{
			Name:       name,
			Identifier: in,

			WarnThreshold: WarnThreshold,
			CritThreshold: CritThreshold,
		}

		order[i] = name
	}

	sort.Strings(order)

	gen := i3.NewOrderedMultiGenerator(batteries, order)

	return gen, nil
}
