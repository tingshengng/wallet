package mock

import (
	"time"
	"wallet/internal/cache"
)

// MockCache is a mock implementation of Cache interface
type MockCache struct {
	cache.Cache
	DeleteFunc func(key string)
	GetFunc    func(key string) (interface{}, bool)
	SetFunc    func(key string, value interface{}, expiration time.Duration)
}

func (m *MockCache) Delete(key string) {
	if m.DeleteFunc != nil {
		m.DeleteFunc(key)
	}
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return nil, false
}

func (m *MockCache) Set(key string, value interface{}, expiration time.Duration) {
	if m.SetFunc != nil {
		m.SetFunc(key, value, expiration)
	}
}
