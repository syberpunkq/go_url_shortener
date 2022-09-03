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
	Key    string `json:"key"`
	Value  string `json:"value"`
	UserId string `json:"user_id"`
}

type Storage struct {
	dict            []map[string]string
	currentIndex    int
	mutex           *sync.RWMutex
	fileStoragePath string
}

func NewStorage() *Storage {
	return &Storage{
		dict:            make([]map[string]string, 0),
		currentIndex:    1,
		mutex:           &sync.RWMutex{},
		fileStoragePath: "",
	}
}

func FileStorage(fileStoragePath string) (*Storage, error) {
	stor := NewStorage()
	file, err := os.OpenFile(fileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(file)

	for {
		var e Entity
		if err := dec.Decode(&e); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		stor.dict = append(stor.dict, map[string]string{
			"key":     e.Key,
			"value":   e.Value,
			"user_id": e.UserId,
		})
	}

	stor.fileStoragePath = fileStoragePath
	file.Close()

	return stor, nil
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
func (s *Storage) Add(val string) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	index := fmt.Sprint(s.currentIndex)
	s.dict[index] = val
	s.currentIndex++

	if s.fileStoragePath != "" {
		err := s.FileAppend(index, val)
		if err != nil {
			log.Print(err)
			return "", err
		}
	}
	return index, nil
}

func (s *Storage) FileAppend(key string, value string) error {
	file, err := os.OpenFile(s.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(Entity{Key: key, Value: value})
	return nil
}
