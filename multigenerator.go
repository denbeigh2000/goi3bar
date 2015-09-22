package goi3bar

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
