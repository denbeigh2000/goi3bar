package goi3bar

import "sync"

type MultiGenerator struct {
	generators []Generator
}

func (m MultiGenerator) Generate() ([]Output, error) {
	out := make([]Output, 0)

	for _, g := range m.generators {
		output, err := g.Generate()
		if err != nil {
			return nil, err
		}

		for _, o := range output {
			out = append(out, o)
		}
	}

	return out, nil
}

func NewMultiGenerator(g []Generator) MultiGenerator {
	return MultiGenerator{g}
}

type OrderedMultiGenerator struct {
	generators map[string]Generator
	order      []string
}

func NewOrderedMultiGenerator(g map[string]Generator, order []string) Generator {
	if len(order) != len(g) {
		panic("Number of keys must be equal")
	}

	for _, key := range order {
		if _, ok := g[key]; !ok {
			panic("Keys must be the same")
		}
	}

	return &OrderedMultiGenerator{g, order}
}

func (g *OrderedMultiGenerator) Generate() ([]Output, error) {
	out := make([]Output, len(g.order))
	wg := sync.WaitGroup{}

	for i, gen := range g.order {
		wg.Add(1)
		go func(i int, g Generator) {
			output, _ := g.Generate()
			out[i] = output[0]
			wg.Done()
		}(i, g.generators[gen])
	}

	wg.Wait()
	return out, nil
}
