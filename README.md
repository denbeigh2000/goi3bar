## goi3bar

[![GoDoc](https://godoc.org/github.com/denbeigh2000/goi3bar?status.svg)](http://godoc.org/github.com/denbeigh2000/goi3bar)

Finally, a configurable, lightweight and easily extensible replacement for i3status.

Why use this over several other alternatives?
 - **Speed**. This performs better than its' cousins written in interpreted languages
   (python, php, etc)
 - **Fine-grained concurrency**. You can assign individual timings to all plugins,
   allowing you to make expensive calls less frequently (think making a network
   call to retrieve the weather, vs. updating the time).
 - **Simple configuration**. goi3bar is driven by JSON configuration, allowing you
   to easily customise your i3bar. Have you ever tried to use conky?
 - **Simple Extensibility**. Writing new plugins is much simpler than writing
   new functionality for a C-based project like conky or i3blocks. There are
   simple interfaces that let you build your own plugins, and handle JSON
   configuration. Look in the godoc for Producer, Genreator and Builder.

There is a sample configuration file in `cmd/goi3bar/config.json`, which contains
configuration for all plugins and all their options.

Talk is cheap! This powers my own i3bar:
![i3bar](http://i.imgur.com/5qwymic.png)

Currently have:
 - Configuration via JSON
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
 - Support for more batteries(?) This was written for a ThinkPad x240 because that's what I have. Pull requests welcome if some battery functionality does not work on your machine. 
