package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// ReadConfigSet reads the JSON file referenced at path, and returns a ConfigSet
// representing that configuration
func ReadConfigSet(r io.Reader) (cs ConfigSet, err error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	cs = ConfigSet{}
	err = json.Unmarshal(data, &cs)

	return
}
