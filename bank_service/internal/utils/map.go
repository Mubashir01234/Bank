package utils

import (
	"encoding/json"
	"sync"
)

// SafeMap is like a Go map[K]V but is safe for concurrent use.
type SafeMap[K comparable, V any] struct {
	data map[K]V
	mtx  *sync.RWMutex
}

// NewSafeMap creates a new instance of SafeMap.
func NewSafeMap[K comparable, V any]() SafeMap[K, V] {
	return SafeMap[K, V]{data: make(map[K]V), mtx: new(sync.RWMutex)}
}

// Load returns the value stored in the map for a key.
func (m SafeMap[K, V]) Load(key K) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return m.data[key]
}

// Store sets the value for a key.
func (m SafeMap[K, V]) Store(key K, value V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.data[key] = value
}

// Delete deletes the value for a key.
func (m SafeMap[K, V]) Delete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.data, key)
}

// Returns the internal map data.
// This is unsafe for concurrent use.
func (m SafeMap[K, V]) Data() map[K]V {
	return m.data
}

func MapToJSON(data map[string]float64) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
