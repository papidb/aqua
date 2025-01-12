package internal

import (
	"bytes"
	"encoding/json"
)

func ToJSON(v interface{}) []byte {
	raw, _ := json.Marshal(v)
	// log API responses
	if v != nil {
		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, raw); err != nil {
			panic(err)
		}
	}
	return raw
}
