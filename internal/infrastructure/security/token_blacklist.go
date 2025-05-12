package security

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	mu    sync.RWMutex
	store map[string]time.Time
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		store: make(map[string]time.Time),
	}
}

func (t *TokenBlacklist) Add(token string, exp time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.store[token] = exp
}

func (t *TokenBlacklist) IsBlacklisted(token string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	exp, exists := t.store[token]
	if !exists {
		return false
	}
	if time.Now().After(exp) {
		delete(t.store, token)
		return false
	}
	return true
}
