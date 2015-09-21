package main

import (
	"bitbucket.org/denbeigh2000/goi3status"
	"time"
)

func formatTime() string {
	time := time.Now()

	return timeFormat.Format(formatString, time)
}

func timeOutput(o *Output) error {
	o.FullText = formatTime()
	return nil
}

func main() {
	clock := NewItem("time", 1*time.Second, timeOutput)

	bar := NewI3bar(1 * time.Second)

	bar.Register("time", clock)

	bar.Start()
	defer bar.Kill()

	<-time.After(20 * time.Second)
}
