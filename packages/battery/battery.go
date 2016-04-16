package battery

import (
	i3 "github.com/denbeigh2000/goi3bar"

	cfg "github.com/alyu/configparser"

	"fmt"
	"strconv"
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

	config, err := cfg.Read(file)
	if err != nil {
		return err
	}

	c, err := config.Section("global")
	if err != nil {
		return err
	}

	b.present = c.ValueOf("POWER_SUPPLY_PRESENT") == "1"
	b.status = c.ValueOf("POWER_SUPPLY_STATUS")

	fullCharge, _ := strconv.ParseFloat(c.ValueOf("POWER_SUPPLY_ENERGY_FULL"), 32)
	currentCharge, _ := strconv.ParseFloat(c.ValueOf("POWER_SUPPLY_ENERGY_NOW"), 32)
	powerUse, _ := strconv.ParseFloat(c.ValueOf("POWER_SUPPLY_POWER_NOW"), 32)

	currentPerc, _ := strconv.Atoi(c.ValueOf("POWER_SUPPLY_CAPACITY"))

	b.level = currentPerc

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
