package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Entity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Storage struct {
	dict            map[string]string
	currentIndex    int
	mutex           *sync.RWMutex
	fileStoragePath string
}

func NewStorage() *Storage {
	return &Storage{
		dict:            make(map[string]string),
		currentIndex:    1,
		mutex:           &sync.RWMutex{},
		fileStoragePath: "",
	}
}

func FileStorage(fileStoragePath string) *Storage {
	stor := NewStorage()
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(file)

	for {
		var e Entity
		if err := dec.Decode(&e); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		stor.dict[e.Key] = e.Value
	}

	stor.fileStoragePath = fileStoragePath
	file.Close()

	return stor
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

	if s.fileStoragePath != "" {
		s.FileAppend(index, val)
	}
	return index
}

func (s *Storage) FileAppend(key string, value string) {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(Entity{Key: key, Value: value})
}
