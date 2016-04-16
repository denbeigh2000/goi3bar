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
   new functionality for a C-based project like conky. There are
   simple interfaces that let you build your own plugins, and handle JSON
   configuration. Look in the godoc for Producer, Genreator and Builder.

Talk is cheap! This powers my own i3bar:

![i3bar1](http://i.imgur.com/3zHCdjv.png)
![i3bar2](http://i.imgur.com/HOTvNyp.png)
![i3bar3](http://i.imgur.com/SnHTnkA.png)

### Getting started

**Dependency**: iwconfig for WLAN info. Should be available in `$PATH`

Either generate a binary with `go build`, or run `go install` in the root dir
and add `$GOPATH/bin` to your `$PATH`, then run the `goi3bar` binary:

From a file:
```
goi3bar --config-path ~/.i3/config.json
```

or through stdin:
```
cat ~/.i3/config.json | goi3bar
```
**NB**: While input through stdin is supported in this version, it will be
removed in a future revision - I am working on support for click events in a
feature branch (it's almost ready - check it out!) which needs to read from
stdin, and we can't read from stdin after reading EOF.

Set either of these as your `status_command` field in `~/.i3/config`

### Basic configuration

A configuration file is represented with JSON, consisting of refresh interval
and zero or more entries

Each entry has a "package" referring to the plugin it uses, a "name" (anything,
but must be unique) and an "options" struct, which will be dependent on the
package you are using.

A set of packages come pre-included in the default "goi3bar" binary

| Package key | Function |
| --- | --- |
| cpu_load | 1, 5, 15 minute CPU loads |
| cpu_util | Current CPU percentage utilisation |
| memory | Current memory usage |
| disk_usage | Current free disk space |
| disk_access | Current data I/O rate |
| battery | Current battery level/remaining time |
| network | Information about currently connected networks |
| clock | Current time |

#### Sample config

This is a very simple config file, there is a sample configuration file in
`cmd/goi3bar/config.json`, which contains configuration for all plugins and
all their options.


```
{
    "colors": {
        "color_crit": "#FF0000",
        "color_warn": "#FFA500",
        "color_ok": "#00FF00",
        "color_general": "#FFFFFF"
    },
    "interval": "5s",
    "entries": [
        {
            "package": "memory",
            "name": "memory",
            "options": {
                "interval": "10s",
                "warn_threshold": 75,
                "crit_threshold": 85
            }
        }
    ]
```

### TODO

Currently have:
 - Support (but no action) for click events
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
