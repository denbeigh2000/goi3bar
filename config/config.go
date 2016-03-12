package config

import (
	. "github.com/denbeigh2000/goi3bar"

	"gopkg.in/yaml.v2"

	"io/ioutil"
	"os"
)

type Builder interface {
	Build(options interface{}) Generator
}

type Config struct {
	Package string      `yaml:"package"`
	Options interface{} `yaml:"options"`
}

type ConfigSet struct {
	Entries []Config `yaml:"entries"`

	builders map[string]Builder
}

func ReadConfigFile(path string) (I3bar, error) {
	return I3bar{}, nil
}

func readConfigSet(path string) (cs ConfigSet, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, cs)
	if err != nil {
		return
	}

	return cs, nil
}

func (c ConfigSet) Build() (I3bar, error) {
	return I3bar{}, nil
}
