// Copyright 2015 Denbeigh Stevens

/*
goi3bar is a package that is capable of generating output suitable for i3bar.
i3bar output is created by outputting JSON, which is defined in the i3bar
documentation, found at http://i3wm.org/docs/i3bar-protocol.html

A single i3bar output struct is represented as an Output. Any applet wanting
to display information on the bar must create output of type []Output.

Two interfaces have been defined to make creating output easier: Generator and
Producer.

Generator: Simple output generator, where output can be created without pausing
to collect data. Canonical examples are a clock, free disk space and free memory.
A generator must be able to generate output on demand. If you need to
periodically collect data, use a Producer.

Producer: Similar to a generator, but is responsible for managing when it
delivers its updates. A Producer is usually used over a generator when there is
a need to specially manage how your collection function is run, so the timing
control is contained within the Producer. A canonical example is the Disk IO
applet, which needs to collect data for a given duration before it can give
average throughput.

To be registered on the bar, the applet needs to implement the Producer
interface. Generators can easily be made into Producers by using the
BaseProducer struct with a Generator, and giving it an interval. The
BaseProducer will call Generate every interval, and push the resultant
updates down the channel.
*/

package goi3bar
