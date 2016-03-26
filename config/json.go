package config

import (
	"encoding/json"
	"io"
)

// ReadConfigSet reads the JSON file referenced at path, and returns a ConfigSet
// representing that configuration
func ReadConfigSet(r io.Reader) (cs ConfigSet, err error) {
	dec := json.NewDecoder(r)
	cs = ConfigSet{}
	err = dec.Decode(&cs)

	return
}
