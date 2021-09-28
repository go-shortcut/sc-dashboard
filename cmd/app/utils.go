package main

import "strconv"

func GetKeysAsInt64Slice(m map[string]interface{}) []int64 {
	keys := make([]int64, len(m))

	i := 0
	for k := range m {
		kInt64, e := strconv.ParseInt(k, 10, 64)
		if e != nil {
			continue
		}
		keys[i] = kInt64
		i++
	}
	return keys
}
