package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"

	"td27/rpc/basis/types/common_pb"
	"td27/rpc/basis/types/sysTool/cache_pb"
)

// blocklistKeyPrefix namespaces blocklist entries in the remote cache store.
// Value is the token expiry as unix seconds.
const blocklistKeyPrefix = "jwt_blocklist:"

const blocklistPageSize = 500

type TokenBlocklist interface {
	IsBlocklisted(tokenID string) bool
	AddToBlocklist(tokenID string, expiresAt time.Time) error
}

type inMemoryBlocklist struct {
	mu      sync.RWMutex
	entries map[string]time.Time
	client  cache_pb.CacheClient // optional remote store for persistence across restarts
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

// InitBlocklistPersistence attaches a remote cache store to the blocklist so
// revoked tokens survive gateway restarts and are shared across instances.
// Must be called before the gateway starts serving traffic.
func InitBlocklistPersistence(client cache_pb.CacheClient) {
	bl := GetBlocklist().(*inMemoryBlocklist)
	bl.setRemoteClient(client)
	bl.refreshFromRemote()
}

func (b *inMemoryBlocklist) setRemoteClient(client cache_pb.CacheClient) {
	b.mu.Lock()
	b.client = client
	b.mu.Unlock()
}

func (b *inMemoryBlocklist) remoteClient() cache_pb.CacheClient {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.client
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// RevokeRequestToken blocklists the JWT presented on the request until its
// natural expiry (or 24h if unparseable). Use after security-sensitive actions
// such as password changes to force a re-login.
func RevokeRequestToken(r *http.Request) {
	tokenStr := strings.TrimPrefix(r.Header.Get("x-token"), "Bearer ")
	if tokenStr == "" {
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok {
				expiresAt = time.Unix(int64(exp), 0)
			}
		}
	}

	if err := GetBlocklist().AddToBlocklist(HashToken(tokenStr), expiresAt); err != nil {
		logx.Errorf("revoke request token failed: %v", err)
	}
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
	b.entries[tokenID] = expiresAt
	b.mu.Unlock()

	if client := b.remoteClient(); client != nil {
		go b.persistToRemote(client, tokenID, expiresAt)
	}
	return nil
}

func (b *inMemoryBlocklist) persistToRemote(client cache_pb.CacheClient, tokenID string, expiresAt time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ttl := int64(time.Until(expiresAt).Seconds())
	if ttl < 1 {
		return
	}
	req := &cache_pb.SetCacheReq{
		Key:        blocklistKeyPrefix + tokenID,
		Value:      strconv.FormatInt(expiresAt.Unix(), 10),
		TtlSeconds: ttl,
	}
	if _, err := b.client.SetCache(ctx, req); err != nil {
		logx.Errorf("persist blocklist entry failed: %v", err)
	}
}

// refreshFromRemote pages through the remote cache store and merges all
// non-expired blocklist entries into the in-memory map.
func (b *inMemoryBlocklist) refreshFromRemote() {
	client := b.remoteClient()
	if client == nil {
		return
	}

	remote := make(map[string]time.Time)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for page := uint32(1); ; page++ {
		resp, err := client.ListCache(ctx, &common_pb.PageReq{Page: page, PageSize: blocklistPageSize})
		if err != nil {
			logx.Errorf("load blocklist from remote failed: %v", err)
			return
		}
		for _, item := range resp.List {
			if len(item.Key) <= len(blocklistKeyPrefix) || item.Key[:len(blocklistKeyPrefix)] != blocklistKeyPrefix {
				continue
			}
			sec, err := strconv.ParseInt(item.Value, 10, 64)
			if err != nil || sec <= 0 {
				continue
			}
			expiresAt := time.Unix(sec, 0)
			if expiresAt.After(time.Now()) {
				remote[item.Key[len(blocklistKeyPrefix):]] = expiresAt
			}
		}
		if len(resp.List) < blocklistPageSize {
			break
		}
	}

	if len(remote) == 0 {
		return
	}

	b.mu.Lock()
	for id, expiresAt := range remote {
		if cur, ok := b.entries[id]; !ok || expiresAt.After(cur) {
			b.entries[id] = expiresAt
		}
	}
	b.mu.Unlock()
}

func (b *inMemoryBlocklist) cleanupLoop() {
	refreshTicker := time.NewTicker(time.Minute)
	cleanupTicker := time.NewTicker(30 * time.Minute)
	defer refreshTicker.Stop()
	defer cleanupTicker.Stop()

	for {
		select {
		case <-refreshTicker.C:
			b.refreshFromRemote()
		case <-cleanupTicker.C:
			b.cleanup()
		}
	}
}

func (b *inMemoryBlocklist) cleanup() {
	b.mu.Lock()
	now := time.Now()
	for id, expiresAt := range b.entries {
		if now.After(expiresAt) {
			delete(b.entries, id)
		}
	}
	b.mu.Unlock()
}
