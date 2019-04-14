package cluster_test

import (
	"testing"
	"time"

	"github.com/alexxeis/keyval/cluster"
	"github.com/alexxeis/keyval/storage"
)

func TestCluster_Get(t *testing.T) {
	c := cluster.NewCluster(10, 0)

	// test missing key
	v, err := c.Get("missing")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test wrong type
	if err = c.Hset("hset", "f", "v"); err != nil {
		t.Error(err)
	}
	if _, err = c.Get("hset"); err != storage.ErrorWrongType {
		t.Error(err)
	}
}

func TestCluster_Set(t *testing.T) {
	c := cluster.NewCluster(10, 0)
	key := "k"

	// test write
	c.Set(key, "v1", 0)
	v, err := c.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "v1" {
		t.Error("expected value = v1, got ", v)
	}

	// test rewrite
	c.Set(key, "v2", 0)
	v, err = c.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "v2" {
		t.Error("expected value = v2, got ", v)
	}

	// test expire
	c.Set(key, "expired", time.Nanosecond)
	time.Sleep(time.Nanosecond)
	v, err = c.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}
}

func TestCluster_Expire(t *testing.T) {
	c := cluster.NewCluster(10, 0)
	key := "k"

	c.Set(key, "v", 0)
	c.Expire(key, time.Nanosecond)
	time.Sleep(time.Nanosecond)

	val, err := c.Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != "" {
		t.Error("not empty value")
	}
}

func TestCluster_Remove(t *testing.T) {
	c := cluster.NewCluster(10, 0)
	key := "k"

	c.Set(key, "v", 0)
	c.Remove(key)

	v, err := c.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}
}

func TestCluster_Keys(t *testing.T) {
	c := cluster.NewCluster(10, 0)

	c.Set("k1", "v", 0)
	c.Set("k2", "v", time.Nanosecond)
	c.Set("k3", "v", 0)
	time.Sleep(time.Nanosecond)

	keys := c.Keys()
	l := len(keys)
	if l != 2 {
		t.Errorf("expected len is %d, got %d", 2, l)
	}

	if keys[0] != "k1" && keys[1] != "k1" {
		t.Error("missing key k1")
	}
	if keys[0] != "k3" && keys[1] != "k3" {
		t.Error("missing key k3")
	}
}

func TestCluster_Hget(t *testing.T) {
	c := cluster.NewCluster(10, 0)

	// test missing key
	v, err := c.Hget("missing", "f")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test missing field
	if err = c.Hset("hset", "f", "v"); err != nil {
		t.Error(err)
	}
	v, err = c.Hget("hset", "missing")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test wrong type
	c.Set("string", "v", 0)
	if _, err = c.Hget("string", "f"); err != storage.ErrorWrongType {
		t.Error(err)
	}
}

func TestCluster_Hset(t *testing.T) {
	c := cluster.NewCluster(10, 0)
	key := "k"
	field1 := "f1"
	field2 := "f2"
	val1 := "v1"
	val2 := "v2"
	val3 := "v3"

	var err error

	// test write
	if err = c.Hset(key, field1, val1); err != nil {
		t.Error(err)
	}
	if err = c.Hset(key, field2, val2); err != nil {
		t.Error(err)
	}

	v, err := c.Hget(key, field1)
	if err != nil {
		t.Error(err)
	}
	if v != val1 {
		t.Errorf("expected value = %s, got %s", val1, v)
	}

	v, err = c.Hget(key, field2)
	if err != nil {
		t.Error(err)
	}
	if v != val2 {
		t.Errorf("expected value = %s, got %s", val2, v)
	}

	// test rewrite
	if err = c.Hset(key, field1, val3); err != nil {
		t.Error(err)
	}

	v, err = c.Hget(key, field1)
	if err != nil {
		t.Error(err)
	}
	if v != val3 {
		t.Errorf("expected value = %s, got %s", val3, v)
	}
}

func TestCluster_Hdel(t *testing.T) {
	c := cluster.NewCluster(10, 0)
	var err error

	// test missing key
	if err = c.Hdel("missing", "missing"); err != nil {
		t.Error(err)
	}

	// test remove
	if err = c.Hset("removed", "removed", "v"); err != nil {
		t.Error(err)
	}
	if err = c.Hdel("removed", "removed"); err != nil {
		t.Error(err)
	}

	v, err := c.Hget("removed", "removed")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test wrong type
	c.Set("string", "v", 0)
	if err = c.Hdel("string", "missing"); err != storage.ErrorWrongType {
		t.Error(err)
	}
}
