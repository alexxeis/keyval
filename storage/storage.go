package storage

import (
	"sync"
	"time"
)

// Storage is the interface for key-value storage
type Storage interface {
	// Shutdown finish storage work
	Shutdown()

	// Expire sets key expiration time
	Expire(key string, ttl time.Duration) bool

	// Set adds value to the storage with the given key that's expire after ttl
	Set(key, val string, ttl time.Duration)

	// Get returns value from the storage by the key
	Get(key string) (string, error)

	// Remove deletes item from the storage by the key
	Remove(key string)

	// Keys returns all key names from the storage
	Keys() []string

	// Hget returns value by key and field
	Hget(key, field string) (string, error)

	// Hset adds value for key and field
	Hset(key, field, val string) error

	// Hdel deletes value by key and field
	Hdel(key, field string) error
}

// item is a basic storage element with data
type item struct {
	value      interface{}
	expiration int64
}

// expired returns true if item is expired
func (i item) expired() bool {
	return i.expiration != 0 && time.Now().UnixNano() > i.expiration
}

// storage is a data storage instance
type storage struct {
	items         map[string]item
	mu            sync.RWMutex
	cleanInterval time.Duration
	done          chan interface{}
}

// NewStorage returns new storage instance
func NewStorage(cleanInterval time.Duration) *storage {
	if cleanInterval < 0 {
		panic("non-positive clean interval")
	}

	s := &storage{
		items:         make(map[string]item),
		cleanInterval: cleanInterval,
		done:          make(chan interface{}),
	}

	if cleanInterval > 0 {
		go s.runCleaner()
	}

	return s
}

// Shutdown stops storage's cleaner
func (s *storage) Shutdown() {
	close(s.done)
}

// runCleaner starts cleaner work
func (s *storage) runCleaner() {
	ticker := time.NewTicker(s.cleanInterval)
	for {
		select {
		case <-ticker.C:
			s.deleteExpiredItems()
		case <-s.done:
			ticker.Stop()
			return
		}
	}
}

// deleteExpiredItems delete all expired items
func (s *storage) deleteExpiredItems() {
	s.mu.Lock()

	for k, v := range s.items {
		if v.expired() {
			delete(s.items, k)
		}
	}

	s.mu.Unlock()
}
