package battery

import (
	i3 "github.com/denbeigh2000/goi3bar"

	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	BaseBatteryPath = "/sys/class/power_supply"
	BatteryPath     = BaseBatteryPath + "/%v/uevent"
)

func calcChargeTime(fullCharge float64, currentCharge float64, powerUse float64) time.Duration {
	return time.Duration(60*60*(fullCharge-currentCharge)/powerUse) * time.Second
}

func calcDepleteTime(currentCharge float64, powerUse float64) time.Duration {
	return time.Duration(60*60*(currentCharge/powerUse)) * time.Second
}

type Battery struct {
	// A friently name used to show the battery in the taskbar. Can be empty.
	Name string
	// A short string used to identify the battery to the system, e.g. BAT0.
	// Should be found in /sys/class/power_supply/<BATTERY_ID>
	Identifier string

	// If the battery level is below this percentage amount, the battery will
	// be rendered int a warning colour
	WarnThreshold int
	// If the battery level is below this percentage amount, the battery will
	// be rendered int a critical colour
	CritThreshold int

	level     int
	present   bool
	status    string
	remaining time.Duration
}

func (b *Battery) Crit() bool {
	return b.level < b.CritThreshold ||
		!b.present
}

func (b *Battery) Warn() bool {
	return b.status == "Discharging" && b.level < b.WarnThreshold
}

func (b *Battery) update() error {
	file := fmt.Sprintf(BatteryPath, b.Identifier)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	var fullCharge, currentCharge, powerUse, currentPerc float64
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tokens := strings.SplitN(scanner.Text(), "=", 2)
		if len(tokens) != 2 {
			continue
		}
		switch tokens[0] {
		case "POWER_SUPPLY_ENERGY_FULL_DESIGN":
			fullCharge, _ = strconv.ParseFloat(tokens[1], 32)
		case "POWER_SUPPLY_CHARGE_FULL":
			fullCharge, _ = strconv.ParseFloat(tokens[1], 32)
		case "POWER_SUPPLY_ENERGY_NOW":
			currentCharge, _ = strconv.ParseFloat(tokens[1], 32)
		case "POWER_SUPPLY_CHARGE_NOW":
			currentCharge, _ = strconv.ParseFloat(tokens[1], 32)
		case "POWER_SUPPLY_STATUS":
			b.status = tokens[1]
		case "POWER_SUPPLY_POWER_NOW":
			powerUse, _ = strconv.ParseFloat(tokens[1], 32)
		case "POWER_SUPPLY_CAPACITY":
			currentPerc, _ = strconv.ParseFloat(tokens[1], 32)
		}
	}

	if powerUse == 0 {
		b.remaining = 0 * time.Hour
	} else {
		switch b.status {
		case "Charging":
			b.remaining = calcChargeTime(fullCharge, currentCharge, powerUse)
		case "Discharging":
			b.remaining = calcDepleteTime(currentCharge, powerUse)
		case "Full":
			b.remaining = 0 * time.Hour
		}
	}

	b.level = int(currentPerc)

	return nil
}

// Generate implements Generator
func (b *Battery) Generate() (out []i3.Output, err error) {
	err = b.update()
	if err != nil {
		return
	}

	o := i3.Output{
		Name:      b.Name,
		Instance:  b.Identifier,
		Separator: true,
	}

	out = make([]i3.Output, 1)
	defer func() {
		out[0] = o
	}()

	if !b.present {
		o.FullText = fmt.Sprintf("Battery %v not present", b.Identifier)
		o.Color = i3.DefaultColors.Crit
		return
	}

	var remain interface{}
	switch b.remaining {
	case time.Duration(0):
		remain = "N/A"
	default:
		remain = b.remaining
	}

	text := fmt.Sprintf("%v %v %v%% %v", b.Name, b.status, b.level, remain)
	switch {
	case b.Crit():
		o.Color = i3.DefaultColors.Crit
	case b.Warn():
		o.Color = i3.DefaultColors.Warn
	default:
		o.Color = i3.DefaultColors.OK
	}

	o.FullText = text

	return
}
