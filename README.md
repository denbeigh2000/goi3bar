## goi3bar

This is a small JSON generator for i3bar, written in golang

I wrote it because I wanted to take advantage of Go's concurrency support to
easily do some tasks less frequently - changing the time every second but
performing an expensive/unimportant operation like, say, checking the weather
less frequently.

Some sample configurations are in the `scripts/` directory, run them with `go run
scripts/scriptname`

I've tried to include some useful interfaces to make making blocks easy, which
I will document in godoc later.

Real men put their money where their mouth is: it powers my own i3bar:

![i3bar](https://i.imgur.com/k1zTMvK.png)
