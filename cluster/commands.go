package cluster

import (
	"sync"
	"time"

	"github.com/alexxeis/keyval/storage"
)

func (c *cluster) Expire(key string, ttl time.Duration) bool {
	return c.instance(key).Expire(key, ttl)
}

func (c *cluster) Set(key, val string, ttl time.Duration) {
	c.instance(key).Set(key, val, ttl)
}

func (c *cluster) Get(key string) (string, error) {
	return c.instance(key).Get(key)
}

func (c *cluster) Remove(key string) {
	c.instance(key).Remove(key)
}

func (c *cluster) Keys() []string {
	ch := make(chan []string, c.count)
	wg := sync.WaitGroup{}
	wg.Add(int(c.count))

	for _, i := range c.instances {
		go func(s storage.Storage, ch chan<- []string, wg *sync.WaitGroup) {
			ch <- s.Keys()
			wg.Done()
		}(i, ch, &wg)
	}

	wg.Wait()
	close(ch)

	var keys []string
	for k := range ch {
		keys = append(keys, k...)
	}

	return keys
}

func (c *cluster) Hget(key, field string) (string, error) {
	return c.instance(key).Hget(key, field)
}

func (c *cluster) Hset(key, field, val string) error {
	return c.instance(key).Hset(key, field, val)
}

func (c *cluster) Hdel(key, field string) error {
	return c.instance(key).Hdel(key, field)
}

func (c *cluster) Shutdown() {
	for _, i := range c.instances {
		i.Shutdown()
	}
}
