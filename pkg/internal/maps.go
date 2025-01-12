package internal

import (
	"strings"
)

// ToLowerKeys converts the keys of a map[string] []string payload into all lower
// cased keys, useful for formatting case insensitive keys like request headers
func ToLowerKeys(headers map[string][]string) map[string]interface{} {
	lowerCaseMap := make(map[string]interface{})

	for k, v := range headers {
		lowerKey := strings.ToLower(k)
		if len(v) == 0 {
			lowerCaseMap[lowerKey] = ""
		} else if len(v) == 1 {
			lowerCaseMap[lowerKey] = v[0]
		} else {
			lowerCaseMap[lowerKey] = v
		}
	}

	return lowerCaseMap
}
