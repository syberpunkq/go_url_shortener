package storage

import "fmt"

var dict = make(map[string]string)
var currentIndex = 1

// Search by short url
// Returns long url and true if found
func FindKey(key string) (string, bool) {
	if value, ok := dict[key]; ok {
		return value, true
	} else {
		return "", false
	}
}

// Search by long url
// Returns index and true if found
func FindVal(val string) (string, bool) {
	// Check if url already in dict
	exists := false
	index := "0"
	for k, v := range dict {
		if v == val {
			exists = true
			index = k
			break
		}
	}
	return index, exists
}

// Adds key-value short-long url, returns index of it
func Add(val string) string {
	index := fmt.Sprint(currentIndex)
	dict[index] = val
	currentIndex++
	return index
}
