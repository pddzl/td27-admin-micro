package sysManagement

import (
	"strings"
	"sync"
	"time"
)

var (
	globalCaptchaStore *CaptchaStore
	captchaStoreOnce   sync.Once
)

func GetCaptchaStore() *CaptchaStore {
	captchaStoreOnce.Do(func() {
		globalCaptchaStore = NewCaptchaStore()
	})
	return globalCaptchaStore
}

type captchaEntry struct {
	code      string
	expiresAt time.Time
}

type CaptchaStore struct {
	mu       sync.RWMutex
	store    map[string]*captchaEntry
	cleanupInterval time.Duration
}

func NewCaptchaStore() *CaptchaStore {
	s := &CaptchaStore{
		store:    make(map[string]*captchaEntry),
		cleanupInterval: 10 * time.Minute,
	}
	go s.cleanupLoop()
	return s
}

func (s *CaptchaStore) Set(id, code string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[id] = &captchaEntry{
		code:      code,
		expiresAt: time.Now().Add(ttl),
	}
}

func (s *CaptchaStore) Verify(id, code string, clear bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.store[id]
	if !ok {
		return false
	}

	if time.Now().After(entry.expiresAt) {
		delete(s.store, id)
		return false
	}

	if !strings.EqualFold(entry.code, code) {
		return false
	}

	if clear {
		delete(s.store, id)
	}
	return true
}

func (s *CaptchaStore) cleanupLoop() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, entry := range s.store {
			if now.After(entry.expiresAt) {
				delete(s.store, id)
			}
		}
		s.mu.Unlock()
	}
}
