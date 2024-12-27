package main

import (
	"sync"
)

type redisMap struct {
	data map[string]string
	mu   sync.Mutex
}

func newRedisMap() *redisMap {
	return &redisMap{data: make(map[string]string)}
}

func (m *redisMap) get(key string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.data[key]
}

func (m *redisMap) set(key string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}
