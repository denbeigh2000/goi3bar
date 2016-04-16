// Copyright 2015-2016 Denbeigh Stevens

/*
Package goi3bar is a package that is capable of generating output suitable for i3bar.
i3bar output is created by outputting JSON, which is defined in the i3bar
documentation, found at http://i3wm.org/docs/i3bar-protocol.html

How to put data on the bar

A single i3bar output struct is represented as an Output. Any applet wanting
to display information on the bar must create output of type []Output. The two
interfaces that produce []Output are Producer and Generator.

Registering new plugins

To register a new plugins on the bar, prodvide a unique key and Producer to
I3bar.Register.
Remember to specify an order of applets before starting the bar. This can be
done by providing all keys in the desired order as a slice to the I3bar.Order
function.

Producer and Generator

The Producer is the interface that must be implemented when registering a plugin
with an I3bar. It has one method, Produce(), which returns a channel of
[]Output. Whenever the I3bar produces new output, it will use the most recently
received []Output it received from that Producer.

It is recommended that you implement this interface only if you have a need to
manage your own output scheduling - if you can generate output in a non-blocking
way, implement the Generator method instead.

Examples of required implementation of the Producer interface can be found in
the disk IO or CPU percentage plugins.

If you do not need to manage your own scheduling, you should implement the
Generator interface instead. It has one method, Generate(), that is assumed to
be non-blocking and will be called at regular intervals. To register a generator
with an I3bar, wrap it in a BaseProducer, which will return a Producer with
managed scheduling.

Examples of the Generator+BaseProducer plugins are the clock, battery, disk usage
and network packages

Clicker

The Clicker is the interface that must be implemented to support interaction
through click events. To properly support this, you must also emit a Name value
in your Output structs, which must match the name your Producer was registered
with, otherwise the event will be unable to be routed back to your Producer.

Builder

The Builder is the interface that must be implemented when making your plugin
available to the user through the configuration file interface. Look in the
config package docs for detailed info on how to provide an API through the
config file.

*/
package goi3bar
