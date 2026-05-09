package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"
)

type TokenBlocklist interface {
	IsBlocklisted(tokenID string) bool
	AddToBlocklist(tokenID string, expiresAt time.Time) error
}

type inMemoryBlocklist struct {
	mu      sync.RWMutex
	entries map[string]time.Time
}

var (
	globalBlocklist *inMemoryBlocklist
	blocklistOnce   sync.Once
)

func GetBlocklist() TokenBlocklist {
	blocklistOnce.Do(func() {
		globalBlocklist = &inMemoryBlocklist{
			entries: make(map[string]time.Time),
		}
		go globalBlocklist.cleanupLoop()
	})
	return globalBlocklist
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (b *inMemoryBlocklist) IsBlocklisted(tokenID string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	expiresAt, ok := b.entries[tokenID]
	if !ok {
		return false
	}

	if time.Now().After(expiresAt) {
		return false
	}

	return true
}

func (b *inMemoryBlocklist) AddToBlocklist(tokenID string, expiresAt time.Time) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries[tokenID] = expiresAt
	return nil
}

func (b *inMemoryBlocklist) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		now := time.Now()
		for id, expiresAt := range b.entries {
			if now.After(expiresAt) {
				delete(b.entries, id)
			}
		}
		b.mu.Unlock()
	}
}
