package battery

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"

	cfg "github.com/alyu/configparser"

	"fmt"
	"strconv"
	"time"
)

const BatteryPath = "/sys/class/power_supply/%v/uevent"

type batInfo struct {
	Status   string `toml:"POWER_SUPPLY_STATUS"`
	Present  bool   `toml:"POWER_SUPPLY_PRESENT"`
	Capacity int    `toml:"POWER_SUPPLY_CAPACITY"`

	FullCharge     int `toml:"POWER_SUPPLY_CHARGE_FULL"`
	CurrentCharge  int `toml:"POWER_SUPPLY_CHARGE_NOW"`
	CurrentCurrent int `toml:"POWER_SUPPLY_CURRENT_NOW"`
}

func NewMultiBattery(names map[string]string, update time.Duration) (map[string]*i3.Item, error) {
	items := make(map[string]*i3.Item, len(names))

	for id, name := range names {
		bat := &Battery{
			Name:       name,
			Identifier: id,
		}

		items[id] = i3.NewItem(id, update, bat)
	}

	return items, nil
}

type Battery struct {
	Name       string
	Identifier string
	Level      int
	Status     string
	Present    bool

	Remaining time.Duration
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

	b.Present = c.ValueOf("POWER_SUPPLY_PRESENT") == "1"
	b.Status = c.ValueOf("POWER_SUPPLY_STATUS")

	fullCharge, _ := strconv.ParseFloat(c.ValueOf("POWER_SUPPLY_ENERGY_FULL"), 32)
	currentCharge, _ := strconv.ParseFloat(c.ValueOf("POWER_SUPPLY_ENERGY_FULL"), 32)
	powerUse, _ := strconv.ParseFloat(c.ValueOf("POWER_SUPPLY_POWER_NOW"), 32)

	currentPerc, _ := strconv.Atoi(c.ValueOf("POWER_SUPPLY_CAPACITY"))

	b.Level = currentPerc

	if powerUse == 0 {
		b.Remaining = 0 * time.Hour
	} else {
		switch b.Status {
		case "Charging":
			b.Remaining = time.Duration(60*60*(fullCharge-currentCharge)/powerUse) * time.Second
		case "Discharging":
			b.Remaining = time.Duration(60*60*(currentCharge/powerUse)) * time.Second
		case "Full":
			b.Remaining = 0 * time.Hour
		}
	}

	return nil
}

func (b *Battery) Update(o *i3.Output) error {
	err := b.update()
	if err != nil {
		return err
	}

	if !b.Present {
		o.FullText = fmt.Sprintf("Battery %v not present", b.Identifier)
		o.Color = "#FF0000"
		return nil
	}

	text := fmt.Sprintf("%v %v %v%% %v", b.Name, b.Status, b.Level, b.Remaining)
	switch {
	case !b.Present || b.Level < 15:
		o.Color = "#FF0000"
	case b.Status == "Charging" || b.Status == "Full" || b.Status == "Discharging" && b.Level >= 35:
		o.Color = "#00FF00"
	default:
		o.Color = "#FFA500"
	}

	o.FullText = text

	return nil
}
