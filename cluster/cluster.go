package cluster

import (
	"time"

	"github.com/alexxeis/keyval/storage"
)

// cluster is a Storage with multi instances support
type cluster struct {
	instances []storage.Storage
	count     int
}

// NewCluster returns new cluster instance
func NewCluster(count int, cleanInterval time.Duration) *cluster {
	if count < 1 {
		panic("wrong cluster instances count")
	}

	instances := make([]storage.Storage, count)
	for i := 0; i < count; i++ {
		instances[i] = storage.NewStorage(cleanInterval)
	}

	return &cluster{
		instances: instances,
		count:     count,
	}
}

// instance returns storage instance by key
func (c *cluster) instance(key string) storage.Storage {
	// TODO: can be better :)
	hasher := newDjb32a()
	hasher.Write([]byte(key))
	sum := int(hasher.Sum32())

	return c.instances[sum%c.count]
}
