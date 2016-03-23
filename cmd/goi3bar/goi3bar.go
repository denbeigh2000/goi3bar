package main

import (
	"github.com/denbeigh2000/goi3bar/config"
	// "github.com/denbeigh2000/goi3bar/packages/battery"
	_ "github.com/denbeigh2000/goi3bar/packages/clock"
	// "github.com/denbeigh2000/goi3bar/packages/memory"
	// "github.com/denbeigh2000/goi3bar/packages/network"

	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// TODO: Flags
const path = "./config.json"

func main() {
	var inFile io.Reader
	switch path {
	case "":
		inFile = os.Stdin
	default:
		f, err := os.Open(path)
		if err != nil {
			panic(fmt.Sprintf("Could not open %v: %v", path, err))
		}
		defer f.Close()

		inFile = f
	}

	cs, err := config.ReadConfigSet(inFile)
	if err != nil {
		panic(err)
	}

	bar, err := cs.Build()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	bar.Start()
	defer bar.Kill()

	for s := range c {
		fmt.Printf("Signal %v received, exiting\n", s)
		return
	}
}
