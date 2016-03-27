package main

import (
	"github.com/denbeigh2000/goi3bar/config"
	_ "github.com/denbeigh2000/goi3bar/packages/battery"
	_ "github.com/denbeigh2000/goi3bar/packages/clock"
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

var path = flag.String("config-path", "stdin", "Path to the config file to use")

func main() {
	flag.Parse()

	var inFile io.Reader
	switch *path {
	case "stdin":
		inFile = os.Stdin
	case "":
		panic("Config path explicitly provided as empty, bailing")
	default:
		f, err := os.Open(*path)
		if err != nil {
			panic(fmt.Sprintf("Could not open %v: %v", path, err))
		}

		inFile = f
	}

	cs, err := config.ReadConfigSet(inFile)
	if err != nil {
		panic(err)
	}

	// Not closing this using defer because otherwise the file stays open until
	// the program is terminated.
	if f, ok := inFile.(*os.File); ok {
		f.Close()
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
