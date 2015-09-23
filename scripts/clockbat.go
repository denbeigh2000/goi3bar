package main

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/packages/battery"
	"github.com/denbeigh2000/goi3bar/packages/clock"

	"time"
)

func main() {

	batteries := make(map[string]i3.Generator, 2)
	batteries["bat1"] = &battery.Battery{
		Name:       "ext bat",
		Identifier: "BAT1",
	}
	batteries["bat0"] = &battery.Battery{
		Name:       "int bat",
		Identifier: "BAT0",
	}

	batteryOrder := []string{"bat0", "bat1"}

	// bats := battery.NewMultiBattery(batteries, batteryOrder)
	batGen := i3.NewOrderedMultiGenerator(batteries, batteryOrder)

	batProd := &i3.BaseProducer{
		Generator: batGen,
		Interval:  5 * time.Second,
		Name:      "bat",
	}

	clockGen := clock.NewClock("%a %d-%b-%y %I:%M:%S")

	clockProd := &i3.BaseProducer{
		Generator: clockGen,
		Interval:  1 * time.Second,
		Name:      "time",
	}

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("time", clockProd)
	bar.Register("bat", batProd)

	bar.Order([]string{"bat", "time"})

	bar.Start()
	defer bar.Kill()

	<-time.After(20 * time.Second)
}
