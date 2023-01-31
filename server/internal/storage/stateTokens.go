package storage

import (
	"errors"
	"fmt"
)

type StateTokenStorage struct {
	values map[string]interface{}
}

var stateTokenStorage Storage

func InitStateTokenStorage() {
	stateTokenStorage = &StateTokenStorage{
		values: make(map[string]interface{}),
	}
}

func GetStateTokenStorage() *Storage {
	return &stateTokenStorage
}

func (s *StateTokenStorage) Add(key string, value interface{}) error {
	s.values[key] = value
	// TODO: Maybe do something usefull with error, eg. check if the key is valid te be added (32 > length)
	return nil
}

func (s *StateTokenStorage) Get(key string) (interface{}, error) {
	if value, ok := s.values[key]; ok {

		return value, nil
	}
	return nil, errors.New("key not found")
}

func (s *StateTokenStorage) Remove(key string) error {
	if _, ok := s.values[key]; ok {
		delete(s.values, key)
		return nil
	}
	return errors.New("key not found")
}

func (s *StateTokenStorage) Move(key string, newKey string) error {
	if value, ok := s.values[key]; ok {
		s.values[newKey] = value
		delete(s.values, key)
		return nil
	}
	return errors.New("key not found")
}

func (s *StateTokenStorage) Clear() {
	s.values = make(map[string]interface{})
}

func (s *StateTokenStorage) String() string {
	return fmt.Sprintf("%v", s.values)
}
