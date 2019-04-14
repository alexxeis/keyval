package storage_test

import (
	"testing"
	"time"

	"github.com/alexxeis/keyval/storage"
)

func TestStorage_Get(t *testing.T) {
	s := storage.NewStorage(0)

	// test missing key
	v, err := s.Get("missing")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test wrong type
	if err = s.Hset("hset", "f", "v"); err != nil {
		t.Error(err)
	}
	if _, err = s.Get("hset"); err != storage.ErrorWrongType {
		t.Error(err)
	}
}

func TestStorage_Set(t *testing.T) {
	s := storage.NewStorage(0)
	key := "k"

	// test write
	s.Set(key, "v1", 0)
	v, err := s.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "v1" {
		t.Error("expected value = v1, got ", v)
	}

	// test rewrite
	s.Set(key, "v2", 0)
	v, err = s.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "v2" {
		t.Error("expected value = v2, got ", v)
	}

	// test expire
	s.Set(key, "expired", time.Nanosecond)
	time.Sleep(time.Nanosecond)
	v, err = s.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}
}

func TestStorage_Expire(t *testing.T) {
	s := storage.NewStorage(0)
	key := "k"

	s.Set(key, "v", 0)
	s.Expire(key, time.Nanosecond)
	time.Sleep(time.Nanosecond)

	val, err := s.Get(key)
	if err != nil {
		t.Error(err)
	}
	if val != "" {
		t.Error("not empty value")
	}
}

func TestStorage_Remove(t *testing.T) {
	s := storage.NewStorage(0)
	key := "k"

	s.Set(key, "v", 0)
	s.Remove(key)

	v, err := s.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}
}

func TestStorage_Keys(t *testing.T) {
	s := storage.NewStorage(0)

	s.Set("k1", "v", 0)
	s.Set("k2", "v", time.Nanosecond)
	s.Set("k3", "v", 0)
	time.Sleep(time.Nanosecond)

	keys := s.Keys()
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

func TestStorage_Hget(t *testing.T) {
	s := storage.NewStorage(0)

	// test missing key
	v, err := s.Hget("missing", "f")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test missing field
	if err = s.Hset("hset", "f", "v"); err != nil {
		t.Error(err)
	}
	v, err = s.Hget("hset", "missing")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test wrong type
	s.Set("string", "v", 0)
	if _, err = s.Hget("string", "f"); err != storage.ErrorWrongType {
		t.Error(err)
	}
}

func TestStorage_Hset(t *testing.T) {
	s := storage.NewStorage(0)
	key := "k"
	field1 := "f1"
	field2 := "f2"
	val1 := "v1"
	val2 := "v2"
	val3 := "v3"

	var err error

	// test write
	if err = s.Hset(key, field1, val1); err != nil {
		t.Error(err)
	}
	if err = s.Hset(key, field2, val2); err != nil {
		t.Error(err)
	}

	v, err := s.Hget(key, field1)
	if err != nil {
		t.Error(err)
	}
	if v != val1 {
		t.Errorf("expected value = %s, got %s", val1, v)
	}

	v, err = s.Hget(key, field2)
	if err != nil {
		t.Error(err)
	}
	if v != val2 {
		t.Errorf("expected value = %s, got %s", val2, v)
	}

	// test rewrite
	if err = s.Hset(key, field1, val3); err != nil {
		t.Error(err)
	}

	v, err = s.Hget(key, field1)
	if err != nil {
		t.Error(err)
	}
	if v != val3 {
		t.Errorf("expected value = %s, got %s", val3, v)
	}
}

func TestStorage_Hdel(t *testing.T) {
	s := storage.NewStorage(0)
	var err error

	// test missing key
	if err = s.Hdel("missing", "missing"); err != nil {
		t.Error(err)
	}

	// test remove
	if err = s.Hset("removed", "removed", "v"); err != nil {
		t.Error(err)
	}
	if err = s.Hdel("removed", "removed"); err != nil {
		t.Error(err)
	}

	v, err := s.Hget("removed", "removed")
	if err != nil {
		t.Error(err)
	}
	if v != "" {
		t.Error("not empty value")
	}

	// test wrong type
	s.Set("string", "v", 0)
	if err = s.Hdel("string", "missing"); err != storage.ErrorWrongType {
		t.Error(err)
	}
}
