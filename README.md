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

### Installation

If you have the pleasure of running Arch Linux, you can simply install the
`goi3bar-git` package from the aur. Otherwise, read on:

1. Install the `go` tool (at least version 1.5)
2. Set up your `GOPATH` and `PATH` environment variables as described
    [here](https://golang.org/doc/code.html#GOPATH).
3. Download this repository and it's dependencies with
    `go get github.com/denbeigh2000/goi3bar/...`
4. Install with `go install github.com/denbeigh2000/goi3bar/cmd/goi3bar`

**Required for WLAN**: iwconfig. Should be available in `$PATH`.
Simple installation check:

```
$ which iwconfig>/dev/null && echo "yay" || echo "no"
yay
```

### Usage

Run the `goi3bar` binary with your config file path as an argument:
```
$ goi3bar -config-path /path/to/your/config.json
```

Set this as the `status_command` field in `~/.i3/config`.

If you see `Error: status_command not found or is missing a library dependency (exit 127)` in your i3bar,
it means `goi3bar` is not in your `$PATH`. Either set your `$PATH` in the script that instantiates i3,
such as `xinitc`, or provide the fully-qualified to the `goi3bar` binary, i.e. `/home/denbeigh/dev/go/bin/goi3bar`

### Configuration

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

#### Sample

This sample config defines an i3bar with custom colours, a 5 second refresh
(not poll) interval, and a single memory printout with colour thresholds
which refreshes every 10 seconds. A full sample config file with all options
configured can be found in `cmd/goi3bar/config.json`.

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
