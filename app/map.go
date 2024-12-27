package main

import (
	"fmt"
	"sync"
	"time"
)

type redisMap struct {
	data      map[string]string
	expireKey map[string]*expirePair
	mu        sync.Mutex
}

type expirePair struct {
	written int64
	px      int
}

func newRedisMap() *redisMap {
	return &redisMap{
		data:      make(map[string]string),
		expireKey: make(map[string]*expirePair),
	}
}

func (m *redisMap) get(key string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	exp := m.expireKey[key]
	fmt.Println("exp pair ", exp)
	if exp != nil {
		ct := time.Now().UnixMilli()
		ext := ct - exp.written
		fmt.Println("time to expire time ", ext)
		fmt.Println("current time ", ct)
		if ext > int64(exp.px) {
			fmt.Println("expired ")
			return ""
		}
		fmt.Println("not expired ")
	}
	return m.data[key]
}

func (m *redisMap) set(key string, value string, px int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	if px != 0 {
		ct := time.Now().UnixMilli()
		m.expireKey[key] = &expirePair{px: px, written: ct}
	}
}
