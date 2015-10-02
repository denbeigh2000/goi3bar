package util

import (
	"fmt"
	"strings"
)

type Bytes float64

const (
	B  Bytes = iota
	KB       = 1 << (10 * iota)
	MB
	GB
	TB
)

func ByteFmt(bytes float64) string {
	unit := ""
	value := Bytes(bytes)
	comp := Bytes(bytes)

	switch {
	case comp >= TB:
		unit = "T"
		value = value / TB
	case comp >= GB:
		unit = "G"
		value = value / GB
	case comp >= MB:
		unit = "M"
		value = value / MB
	case comp >= KB:
		unit = "K"
		value = value / KB
	case comp >= B:
		unit = "B"
	case comp == 0:
		return "0B"
	}

	stringValue := fmt.Sprintf("%.1f", value)
	stringValue = strings.TrimSuffix(stringValue, ".0")
	return fmt.Sprintf("%s%s", stringValue, unit)
}
