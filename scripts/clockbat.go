package main

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"
	"bitbucket.org/denbeigh2000/goi3bar/packages/battery"
	"bitbucket.org/denbeigh2000/goi3bar/packages/clock"

	"time"
)

func main() {

	batteries := make(map[string]*battery.Battery, 2)
	batteries["bat1"] = &battery.Battery{
		Name:       "ext bat",
		Identifier: "BAT1",
	}
	batteries["bat0"] = &battery.Battery{
		Name:       "int bat",
		Identifier: "BAT0",
	}

	batteryOrder := []string{"bat0", "bat1"}

	bats := battery.NewMultiBattery(batteries, batteryOrder)

	batProd := &i3.BaseProducer{
		Generator: &bats,
		Interval:  5 * time.Second,
		Name:      "bat",
	}

	clock := clock.NewClock("%a %d-%b-%y %I:%M:%S")

	clockProd := &i3.BaseProducer{
		Generator: clock,
		Interval:  1 * time.Second,
		Name:      "time",
	}

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("time", clockProd)
	bar.Register("bat", batProd)
	// bar.Register("BAT0", bats["BAT0"])
	// bar.Register("BAT1", bats["BAT1"])

	bar.Order([]string{"bat", "time"})

	bar.Start()
	defer bar.Kill()

	<-time.After(20 * time.Second)
}
