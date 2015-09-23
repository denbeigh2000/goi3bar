package main

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/packages/clock"

	"time"
)

func main() {

	clock := clock.Clock{
		Format: "%a %d-%b-%y %I:%M:%S",
	}

	item := &i3.BaseProducer{
		Generator: clock,
		Interval:  1 * time.Second,
		Name:      "time",
	}

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("time", item)

	bar.Start()
	defer bar.Kill()

	select {}
}
