package main

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/packages/battery"
	"github.com/denbeigh2000/goi3bar/packages/clock"
	"github.com/denbeigh2000/goi3bar/packages/cpu"
	"github.com/denbeigh2000/goi3bar/packages/memory"

	"time"
)

func main() {
	cpu := cpu.Cpu{
		WarnThreshold: 0.7,
		CritThreshold: 0.95,
	}
	cpuProd := &i3.BaseProducer{
		Generator: &cpu,
		Interval:  5 * time.Second,
		Name:      "cpu",
	}

	batteryNames := map[string]string{
		"BAT0": "int bat",
		"BAT1": "ext bat",
	}

	batteries, err := battery.BatteryDiscover(batteryNames, 35, 15)
	if err != nil {
		panic(err)
	}

	mem := memory.Memory{}
	memProd := &i3.BaseProducer{
		Generator: mem,
		Interval:  5 * time.Second,
		Name:      "mem",
	}

	batProd := &i3.BaseProducer{
		Generator: batteries,
		Interval:  5 * time.Second,
		Name:      "bat",
	}

	clockGen := clock.Clock{
		Format: "%a %d-%b-%y %I:%M:%S",
		Color:  "#FFFFFF",
	}

	clockProd := &i3.BaseProducer{
		Generator: clockGen,
		Interval:  1 * time.Second,
		Name:      "time",
	}

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("time", clockProd)
	bar.Register("bat", batProd)
	bar.Register("mem", memProd)
	bar.Register("cpu", cpuProd)

	bar.Order([]string{"cpu", "mem", "bat", "time"})

	bar.Start()
	defer bar.Kill()

	select {}
}
