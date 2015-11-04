package main

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/packages/battery"
	"github.com/denbeigh2000/goi3bar/packages/clock"
	"github.com/denbeigh2000/goi3bar/packages/cpu"
	"github.com/denbeigh2000/goi3bar/packages/disk"
	"github.com/denbeigh2000/goi3bar/packages/memory"
	"github.com/denbeigh2000/goi3bar/packages/network"

	"time"
)

func main() {
	home := disk.DiskUsageItem{
		Name: "home",
		Path: "/home",
	}

	root := disk.DiskUsageItem{
		Name: "root",
		Path: "/",
	}

	diskFree := disk.DiskUsageGenerator{
		WarnThreshold: 25,
		CritThreshold: 15,

		Items: []disk.DiskUsageItem{root, home},
	}

	diskFreeProd := &i3.BaseProducer{
		Generator: diskFree,
		Interval:  1 * time.Minute,
	}

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

	cpuGen := cpu.Cpu{
		WarnThreshold: 0.7,
		CritThreshold: 0.95,
	}
	cpuProd := &i3.BaseProducer{
		Generator: &cpuGen,
		Interval:  5 * time.Second,
		Name:      "cpu",
	}

	cpuPerc := &cpu.CpuPerc{
		WarnThreshold: 75.0,
		CritThreshold: 90.0,
		Interval:      2 * time.Second,
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

	batGen := i3.NewOrderedMultiGenerator(batteries, batteryOrder)

	batProd := &i3.BaseProducer{
		Generator: batGen,
		Interval:  5 * time.Second,
		Name:      "bat",
	}

	SFClockGen := clock.Clock{
		Format:   "SF: %a %d-%b-%y %I:%M:%S",
		Location: "America/Los_Angeles",
	}

	SydClockGen := clock.Clock{
		Format:   "Syd: %a %d-%b-%y %I:%M:%S",
		Location: "Australia/Sydney",
	}

	SFClockProd := &i3.BaseProducer{
		Generator: SFClockGen,
		Interval:  1 * time.Second,
		Name:      "SFTime",
	}

	SydClockProd := &i3.BaseProducer{
		Generator: SydClockGen,
		Interval:  1 * time.Second,
		Name:      "SydTime",
	}

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("SydTime", SydClockProd)
	bar.Register("SFTime", SFClockProd)
	bar.Register("bat", batProd)
	bar.Register("mem", memProd)
	bar.Register("cpu", cpuProd)
	bar.Register("cpuPerc", cpuPerc)
	bar.Register("net", netProd)
	bar.Register("disk", diskFreeProd)

	bar.Order([]string{"cpu", "cpuPerc", "mem", "disk", "bat", "net", "SydTime", "SFTime"})

	bar.Start()
	defer bar.Kill()

	select {}
}
