package main

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/packages/battery"
	"github.com/denbeigh2000/goi3bar/packages/clock"
	"github.com/denbeigh2000/goi3bar/packages/cpu"
	"github.com/denbeigh2000/goi3bar/packages/memory"
	"github.com/denbeigh2000/goi3bar/packages/network"

	"time"
)

func main() {
	wlan := network.BasicNetworkDevice{
		Name:       "wifi",
		Identifier: "wlp3s0",
	}

	wlanDevice := network.WLANDevice{
		BasicNetworkDevice: wlan,
	}

	eth := network.BasicNetworkDevice{
		Name:       "ethernet",
		Identifier: "enp0s25",
	}

	net := network.MultiDevice{
		Devices: map[string]network.NetworkDevice{
			"wlp3s0":  &wlanDevice,
			"enp0s25": &eth,
		},
		Preference: []string{"wlp3s0", "enp0s25"},
	}

	netProd := &i3.BaseProducer{
		Generator: &net,
		Interval:  2 * time.Second,
		Name:      "net",
	}

	cpu := cpu.Cpu{
		WarnThreshold: 0.7,
		CritThreshold: 0.95,
	}
	cpuProd := &i3.BaseProducer{
		Generator: &cpu,
		Interval:  5 * time.Second,
		Name:      "cpu",
	}
	mem := memory.Memory{}
	batteries := make(map[string]i3.Generator, 2)
	batteries["bat1"] = &battery.Battery{
		Name:       "ext bat",
		Identifier: "BAT1",
	}

	memProd := &i3.BaseProducer{
		Generator: mem,
		Interval:  5 * time.Second,
		Name:      "mem",
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

	clockGen := clock.Clock{
		Format: "%a %d-%b-%y %I:%M:%S",
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
	bar.Register("net", netProd)

	bar.Order([]string{"cpu", "mem", "bat", "net", "time"})

	bar.Start()
	defer bar.Kill()

	select {}
}
