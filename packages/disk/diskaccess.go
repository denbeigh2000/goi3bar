package disk

import (
	i3 "github.com/denbeigh2000/goi3bar"
	"github.com/denbeigh2000/goi3bar/util"

	psDisk "github.com/shirou/gopsutil/disk"

	"fmt"
	"log"
	"strconv"
	"time"
)

// Returns a map of the current IO byte count for all devices
func currentIOStatus(devices ...string) (map[string]uint64, error) {
	send := make(map[string]uint64, len(devices))

	allUsage, err := psDisk.IOCounters()
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		usage := allUsage[d]
		send[d] = usage.WriteBytes + usage.ReadBytes
	}

	return send, nil
}

// generateIO generates a map of interface name to average disk IO for the
// past interval each interval, and sends it down the channel. Any errors are
// sent down the error channel
func generateIO(kill <-chan struct{}, interval time.Duration,
	devices ...string) (<-chan map[string]float64, chan error) {

	out := make(chan map[string]float64)
	errCh := make(chan error)

	go func() {
		tick := util.NewTicker(interval, true)
		defer tick.Kill()

		secs := interval.Seconds()
		// Keeps a record of our last checked IO count, so we can calculate our delta
		prev := make(map[string]uint64, len(devices))

		for {
			select {
			case <-tick.C:
				send := make(map[string]float64, len(devices))

				stats, err := currentIOStatus(devices...)
				if err != nil {
					errCh <- err
					continue
				}

				for _, d := range devices {
					// Retrieve the IO stats of this device
					usage, ok := stats[d]
					if !ok {
						errCh <- fmt.Errorf("No device found: %v", d)
						continue
					}

					// Retrieve previous value from cache
					pr, ok := prev[d]
					if !ok {
						// If unavailable send zero, then update cache
						prev[d] = usage
						send[d] = 0
						continue
					}

					diff := float64(usage - pr)

					// Create average throughput by dividing by number of secs in interval
					val := diff / secs
					if diff == 0 {
						val = 0.0
					}
					// Update map to send
					send[d] = val
					// Update cache
					prev[d] = usage
				}

				out <- send
			case <-kill:
				return
			}
		}
	}()

	return out, errCh
}

func (g *DiskIOGenerator) createOutput(i DiskIOItem, v float64) i3.Output {
	text := fmt.Sprintf("%v: %v/s", i.Name, util.ByteFmt(v))

	var color string
	switch {
	case v > g.CritThreshold:
		color = i3.DefaultColors.Crit
	case v > g.WarnThreshold:
		color = i3.DefaultColors.Warn
	default:
		color = i3.DefaultColors.OK
	}

	return i3.Output{
		FullText: text,
		Color:    color,
	}
}

// Produce implements i3.Producer
func (g *DiskIOGenerator) Produce(kill <-chan struct{}) <-chan []i3.Output {
	out := make(chan []i3.Output)

	devices := make([]string, len(g.Items))
	deviceMap := make(map[string]DiskIOItem, len(g.Items))
	for i, item := range g.Items {
		devices[i] = item.Device
		deviceMap[item.Device] = item
	}

	in, errCh := generateIO(kill, g.Interval, devices...)

	go func() {
		for {
			select {
			case err := <-errCh:
				log.Printf("Disk IO error: %v", err)
			case <-kill:
				return
			}
		}
	}()

	go func() {
		defer close(out)

		for {
			select {
			case items := <-in:
				output := make([]i3.Output, len(g.Items))
				for i, item := range g.Items {
					val := items[item.Device]
					outPart := g.createOutput(item, val)
					outPart.Name = g.Name
					outPart.Instance = strconv.Itoa(i)
					output[i] = outPart
				}

				if len(output) > 0 {
					// Add a divider between this and the next one
					output[len(output)-1].Separator = true
				}

				out <- output
			case <-kill:
				return
			}
		}
	}()

	return out
}

type DiskIOItem struct {
	Name   string `json:"name"`
	Device string `json:"device"`
}

type DiskIOGenerator struct {
	WarnThreshold float64
	CritThreshold float64

	Name string

	Interval time.Duration

	Items []DiskIOItem
}
