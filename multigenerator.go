package goi3bar

import (
	"log"
	"sync"
)

// MultiGenerator is a Generator that combines the output of multiple
// Generators. Consistent order is not guaranteed across multiple generations
type MultiGenerator struct {
	generators []Generator
}

// NewMultiGenerator takes a slice of generators and returns a single
// Generator
func NewMultiGenerator(g []Generator) MultiGenerator {
	return MultiGenerator{g}
}

// Generate implements Generator
func (m MultiGenerator) Generate() ([]Output, error) {
	var out []Output

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

// OrderedMultiGenerator is a Generator that can generate output for multiple
// Generators, and keeps the order of outputs the same.
type OrderedMultiGenerator struct {
	generators map[string]Generator
	order      []string
}

// NewOrderedMultiGenerator creates a new OrderedMultiGenerator. It takes a
// map of key -> generator pairs, as well as a slice of keys specifying the
// order the items should appear on the bar.
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

// Generate implements Generator
func (g *OrderedMultiGenerator) Generate() ([]Output, error) {
	out := make([]Output, len(g.order))
	wg := sync.WaitGroup{}

	for i, gen := range g.order {
		wg.Add(1)
		go func(i int, g Generator) {
			output, err := g.Generate()
			if err != nil {
				log.Printf("Error concurrently generating: %v\n", err)
			}
			if len(output) > 0 {
				out[i] = output[0]
			}
			wg.Done()
		}(i, g.generators[gen])
	}

	wg.Wait()
	return out, nil
}
