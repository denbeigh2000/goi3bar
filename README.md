## goi3bar

[![GoDoc](https://godoc.org/github.com/denbeigh2000/goi3bar?status.svg)](http://godoc.org/github.com/denbeigh2000/goi3bar)

This is a small JSON generator for i3bar, written in golang

I wrote it because I wanted to take advantage of Go's concurrency support to
easily do some tasks less frequently - changing the time every second but
performing an expensive/unimportant operation like, say, checking the weather
less frequently.

Some sample configurations are in the `scripts/` directory, run them with `go run
scripts/scriptname`

I've tried to include some useful interfaces to make making blocks easy, which
I will document in godoc later.

Real men put their money where their mouth is! This powers my own i3bar:
![i3bar](https://i.imgur.com/iewiNim.png)

Currently have:
 - Formattable clock
 - Memory usage (with configurable color thresholds)
 - CPU load averages (with configurable color thresholds)
 - Battery values (with automagic discovery and configurable thresholds)
 - Network info with funky applet which only shows most preferred connected network
 - Disk read/write rates
 - Disk usage

Want to have:
 - Unit testing!
 - More configurability for memory, battery moinitors (e.g., formattable)
