package main

import (
	"github.com/denbeigh2000/goi3bar/config"
	// "github.com/denbeigh2000/goi3bar/packages/battery"
	_ "github.com/denbeigh2000/goi3bar/packages/clock"
	// "github.com/denbeigh2000/goi3bar/packages/memory"
	// "github.com/denbeigh2000/goi3bar/packages/network"
)

// TODO: Flags
const path = "./config.json"

func main() {
	cs, err := config.ReadConfigSet(path)
	if err != nil {
		panic(err)
	}

	bar, err := cs.Build()
	if err != nil {
		panic(err)
	}

	bar.Start()
	defer bar.Kill()

	select {}
}
