package storage

import (
	"errors"
	"time"
)

var (
	ErrorWrongType = errors.New("operation against a key holding the wrong kind of value")
)

// getExpiration returns expiration timestamp by TTL
func getExpiration(ttl time.Duration) int64 {
	if ttl < 0 {
		panic("ttl cant be less than 0")
	}

	if ttl > 0 {
		return time.Now().Add(ttl).UnixNano()
	}

	return 0
}

func (s *storage) Expire(key string, ttl time.Duration) bool {
	exp := getExpiration(ttl)

	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.items[key]
	if !ok {
		return false
	}

	i.expiration = exp
	s.items[key] = i
	return true
}

func (s *storage) Set(key, val string, ttl time.Duration) {
	exp := getExpiration(ttl)

	s.mu.Lock()
	s.items[key] = item{
		value:      val,
		expiration: exp,
	}
	s.mu.Unlock()
}

func (s *storage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	i, ok := s.items[key]
	if !ok {
		return "", nil
	}

	if i.expired() {
		return "", nil
	}

	val, ok := i.value.(string)
	if !ok {
		return "", ErrorWrongType
	}

	return val, nil
}

func (s *storage) Remove(key string) {
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}

func (s *storage) Keys() []string {
	s.mu.RLock()

	keys := make([]string, 0, len(s.items))
	for k, v := range s.items {
		if !v.expired() {
			keys = append(keys, k)
		}
	}

	s.mu.RUnlock()
	return keys
}

func (s *storage) Hget(key, field string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	i, ok := s.items[key]
	if !ok {
		return "", nil
	}

	if i.expired() {
		return "", nil
	}

	hmap, ok := i.value.(map[string]string)
	if !ok {
		return "", ErrorWrongType
	}

	return hmap[field], nil
}

func (s *storage) Hset(key, field, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.items[key]
	if !ok {
		s.items[key] = item{
			value: map[string]string{
				field: val,
			},
		}
		return nil
	}

	hmap, ok := i.value.(map[string]string)
	if !ok {
		return ErrorWrongType
	}

	hmap[field] = val
	return nil
}

func (s *storage) Hdel(key, field string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, ok := s.items[key]
	if !ok {
		return nil
	}

	hmap, ok := i.value.(map[string]string)
	if !ok {
		return ErrorWrongType
	}

	delete(hmap, field)
	return nil
}
