package utils

import "sort"

type AttrsMap map[string]string

func (a AttrsMap) Keys() []string {
	keys := make([]string, 0)
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (a AttrsMap) ContainsKey(value string) bool {
	for k := range a {
		if k == value {
			return true
		}
	}
	return false
}
