package util

import "encoding/json"

// This feels so wrong and it really oughtn't to exist. However, this really
// appears to be the commonly accepted way to "defer" decoding parts of a JSON
// blob, so to speak. I have tried all forms of modification and casting, have
// looked into fatih's excellent structs package, tried recursively descending
// the tree by casting it to a map[string]interface{} and encoding as I go, but
// this is the way it has to be.
func JSONReparse(i, o interface{}) (err error) {
	marshalled, err := json.Marshal(i)
	if err == nil {
		err = json.Unmarshal(marshalled, o)
	}

	return
}
