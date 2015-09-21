package main

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"
	"bitbucket.org/denbeigh2000/goi3bar/packages"

	"time"
)

func main() {

	clock := clock.NewClock("%a %d-%b-%y %I:%M:%S")

	item := i3.NewItem("time", 1*time.Second, clock)

	bar := i3.NewI3bar(1 * time.Second)

	bar.Register("time", item)

	bar.Start()
	defer bar.Kill()

	<-time.After(20 * time.Second)
}
