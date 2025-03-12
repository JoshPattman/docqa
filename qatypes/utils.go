package qatypes

import (
	"fmt"
	"strings"
)

func cleanName(name string) string {
	return strings.Trim(name, " \n\r\t")
}

func get[T any](m map[string]any, k string) (T, error) {
	if val, ok := m[k]; !ok {
		return *new(T), fmt.Errorf("failed to find key %s in map %s", k, m)
	} else {
		if typedVal, ok := val.(T); !ok {
			return *new(T), fmt.Errorf("could not cast val %s to type %T", val, typedVal)
		} else {
			return typedVal, nil
		}
	}
}
