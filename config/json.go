package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// This feels so wrong and it really oughtn't to exist. However, this really
// appears to be the commonly accepted way to "defer" decoding parts of a JSON
// blob, so to speak. I have tried all forms of modification and casting, have
// looked into fatih's excellent structs package, tried recursively descending
// the tree by casting it to a map[string]interface{} and encoding as I go, but
// this is the way it has to be.
func jsonReparse(i, o interface{}) (err error) {
	marshalled, err := json.Marshal(i)
	if err == nil {
		err = json.Unmarshal(marshalled, o)
	}

	return
}

// ReadConfigSet reads the JSON file referenced at path, and returns a ConfigSet
// representing that configuration
func ReadConfigSet(path string) (cs ConfigSet, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	cs = ConfigSet{}
	err = json.Unmarshal(data, &cs)
	if err != nil {
		return
	}

	return cs, nil
}
