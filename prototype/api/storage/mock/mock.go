package mock

import (
	"github.com/pbivrell/lights/api/storage"
)

func NewMockStorage[T any]() *MockStorage[T] {
	return &MockStorage[T]{
		store: make(map[string]T),
	}
}

type MockStorage[T any] struct {
	store map[string]T
}

func (m *MockStorage[T]) Write(key string, value *T) error {
	m.store[key] = *value
	return nil
}

func (m *MockStorage[T]) Delete(key string) error {
	delete(m.store, key)
	return nil
}

func (m *MockStorage[T]) Read(key string) (*T, error) {

	value, ok := m.store[key]
	if !ok {
		return nil, storage.ErrorNotFound
	}
	return &value, nil
}

func (m *MockStorage[T]) List() ([]*T, error) {
	results := make([]*T, len(m.store))

	idx := 0
	for _, v := range m.store {
		results[idx] = &v
		idx++
	}
	return results, nil
}

func AppStorage() storage.AppStorage {
	return storage.AppStorage{
		User:    NewMockStorage[storage.User](),
		Session: NewMockStorage[storage.Session](),
		Hub:     NewMockStorage[storage.Hub](),
		Light:   NewMockStorage[storage.Light](),
		Pattern: NewMockStorage[storage.Pattern](),
	}
}
