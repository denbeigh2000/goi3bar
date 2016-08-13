package main

import (
	"github.com/denbeigh2000/goi3bar/config"
	_ "github.com/denbeigh2000/goi3bar/packages/battery"
	_ "github.com/denbeigh2000/goi3bar/packages/clock"
	_ "github.com/denbeigh2000/goi3bar/packages/command"
	_ "github.com/denbeigh2000/goi3bar/packages/cpu"
	_ "github.com/denbeigh2000/goi3bar/packages/disk"
	_ "github.com/denbeigh2000/goi3bar/packages/memory"
	_ "github.com/denbeigh2000/goi3bar/packages/network"

	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var path = flag.String("config-path", "", "Path to the config file to use")

var noPathErr = fmt.Errorf("A configuration file is required, specify one with the -config-json flag")

func loadConfigSet() (config.ConfigSet, error) {
	var inFile io.Reader
	switch *path {
	case "":
		return config.ConfigSet{}, noPathErr
	default:
		f, err := os.Open(*path)
		if err != nil {
			return config.ConfigSet{}, fmt.Errorf("Could not open %v: %v", path, err)
		}
		defer f.Close()

		inFile = f
	}

	return config.ReadConfigSet(inFile)
}

func main() {
	flag.Parse()

	cs, err := loadConfigSet()
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

	bar.Start(os.Stdin)
	defer bar.Kill()

	for s := range c {
		fmt.Printf("Signal %v received, exiting\n", s)
		return
	}
}
