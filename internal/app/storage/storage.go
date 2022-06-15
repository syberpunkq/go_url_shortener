package storage

import (
	"fmt"
	"sync"
)

type Storage struct {
	dict         map[string]string
	currentIndex int
	mutex        *sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		dict:         make(map[string]string),
		currentIndex: 1,
		mutex:        &sync.RWMutex{},
	}
}

// Search by short url
// Returns long url and true if found
func (s Storage) FindKey(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if value, ok := s.dict[key]; ok {
		return value, true
	} else {
		return "", false
	}
}

// Search by long url
// Returns index and true if found
func (s Storage) FindVal(val string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Check if url already in dict
	exists := false
	index := "0"
	for k, v := range s.dict {
		if v == val {
			exists = true
			index = k
			break
		}
	}
	return index, exists
}

// Adds key-value short-long url, returns index of it
func (s *Storage) Add(val string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	index := fmt.Sprint(s.currentIndex)
	s.dict[index] = val
	s.currentIndex++
	return index
}
