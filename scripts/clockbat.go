package main

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"
	"bitbucket.org/denbeigh2000/goi3bar/packages/battery"
	"bitbucket.org/denbeigh2000/goi3bar/packages/clock"

	"time"
)

func main() {

	clock := clock.NewClock("%a %d-%b-%y %I:%M:%S")
	batteryNames := make(map[string]string, 2)
	batteryNames["BAT1"] = "ext bat"
	batteryNames["BAT0"] = "int bat"
	bats, err := battery.NewMultiBattery(batteryNames, 5*time.Second)
	if err != nil {
		panic(err)
	}

	item := i3.NewItem("time", 1*time.Second, clock)

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register(item)
	bar.Register(bats["BAT0"])
	bar.Register(bats["BAT1"])

	bar.Order([]string{"bat0", "bat1", "time"})

	bar.Start()
	defer bar.Kill()

	<-time.After(20 * time.Second)
}
