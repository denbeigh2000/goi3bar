package main

import (
	"github.com/denbeigh2000/goi3bar/config"
	"github.com/denbeigh2000/goi3bar/packages/battery"
	"github.com/denbeigh2000/goi3bar/packages/clock"
	"github.com/denbeigh2000/goi3bar/packages/memory"
)

// TODO: Flags
const path = "./config.json"

func main() {
	cs, err := config.ReadConfigSet(path)
	if err != nil {
		panic(err)
	}

	cs.Register("memory", memory.Build)
	cs.Register("clock", clock.Build)
	cs.Register("battery", battery.Build)

	bar, err := cs.Build()
	if err != nil {
		panic(err)
	}

	bar.Start()
	defer bar.Kill()

	select {}
}
