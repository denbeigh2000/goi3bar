package battery

import (
	i3 "bitbucket.org/denbeigh2000/goi3bar"
	"sync"
)

type MultiBattery struct {
	batteries map[string]*Battery
	order     []string
}

func NewMultiBattery(b map[string]*Battery, order []string) MultiBattery {
	return MultiBattery{b, order}
}

func (m MultiBattery) Generate() ([]i3.Output, error) {
	out := make([]i3.Output, len(m.order))
	wg := sync.WaitGroup{}

	for i, bat := range m.order {
		wg.Add(1)
		go func(i int, b *Battery) {
			output, _ := b.Generate()
			out[i] = output[0]
			wg.Done()
		}(i, m.batteries[bat])
	}

	wg.Wait()
	return out, nil
}
