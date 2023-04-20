package cleo

import (
	"os"
	"sync"
)

type Env struct {
	data map[string]string
	mu   sync.RWMutex
}

func (e *Env) Get(key string) string {
	if e == nil {
		return os.Getenv(key)
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.data == nil {
		return os.Getenv(key)
	}

	return e.data[key]
}

func (e *Env) Set(key, value string) {
	if e == nil {
		os.Setenv(key, value)
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if e.data == nil {
		e.data = map[string]string{}
	}

	e.data[key] = value
}
