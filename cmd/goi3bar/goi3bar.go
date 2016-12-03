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
	"log"
	"os"
	"os/signal"
	"syscall"
)

var path = flag.String("config-path", "", "Path to the config file to use")

var noPathErr = fmt.Errorf("A configuration file is required, see `goi3bar -h`")

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
		log.Fatalf("Error parsing configuration file: %v", err)
	}

	bar, err := cs.Build()
	if err != nil {
		log.Fatalf("Error generating configuration: %v", err)
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM)

	bar.Start(os.Stdin)
	defer bar.Kill()

	for s := range c {
		log.Printf("Signal %v received, exiting\n", s)
		return
	}
}
