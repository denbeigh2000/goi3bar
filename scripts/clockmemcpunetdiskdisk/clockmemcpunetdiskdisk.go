package main

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/packages/clock"
	"github.com/denbeigh2000/goi3bar/packages/cpu"
	"github.com/denbeigh2000/goi3bar/packages/disk"
	"github.com/denbeigh2000/goi3bar/packages/memory"
	"github.com/denbeigh2000/goi3bar/packages/network"

	"time"
)

const Megabyte float64 = 1024 * 1024

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

	eth := network.BasicNetworkDevice{
		Name:       "ethernet",
		Identifier: "enp0s25",
	}

	net := network.MultiDevice{
		Devices: map[string]network.NetworkDevice{
			"enp0s25": &eth,
		},
		Preference: []string{"enp0s25"},
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

	memProd := &i3.BaseProducer{
		Generator: mem,
		Interval:  5 * time.Second,
		Name:      "mem",
	}

	clockGen := clock.Clock{
		Format: "%a %d-%b-%y %I:%M:%S",
	}

	clockProd := &i3.BaseProducer{
		Generator: clockGen,
		Interval:  1 * time.Second,
		Name:      "time",
	}

	diskIO := &disk.DiskIOGenerator{
		WarnThreshold: 3 * Megabyte,
		CritThreshold: 10 * Megabyte,

		Interval: 2 * time.Second,

		Items: []disk.DiskIOItem{
			disk.DiskIOItem{
				Name:   "/home",
				Device: "sde2",
			},
			disk.DiskIOItem{
				Name:   "/",
				Device: "sdc5",
			},
		},
	}
	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("time", clockProd)
	bar.Register("mem", memProd)
	bar.Register("cpu", cpuProd)
	bar.Register("net", netProd)
	bar.Register("disk", diskFreeProd)
	bar.Register("diskio", diskIO)

	bar.Order([]string{"diskio", "disk", "cpu", "mem", "net", "time"})

	bar.Start()
	defer bar.Kill()

	select {}
}
