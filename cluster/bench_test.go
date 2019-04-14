package cluster_test

import (
	"sync"
	"testing"

	"github.com/alexxeis/keyval/cluster"
	"github.com/alexxeis/keyval/storage"
)

func setConcurrently(i int, s storage.Storage, wg *sync.WaitGroup) {
	s.Set(string(i), "v", 0)
	wg.Done()
}

func getConcurrently(i int, s storage.Storage, wg *sync.WaitGroup) {
	s.Get("k" + string(i))
	wg.Done()
}

func BenchmarkStorage_Set(b *testing.B) {
	s := storage.NewStorage(0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go setConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkCluster_Set10(b *testing.B) {
	s := cluster.NewCluster(10, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go setConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkCluster_Set100(b *testing.B) {
	s := cluster.NewCluster(100, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go setConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkCluster_Set1000(b *testing.B) {
	s := cluster.NewCluster(1000, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go setConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkStorage_SetGet(b *testing.B) {
	s := storage.NewStorage(0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go setConcurrently(i, s, &wg)
		go getConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkCluster_SetGet10(b *testing.B) {
	s := cluster.NewCluster(10, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go setConcurrently(i, s, &wg)
		go getConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkCluster_SetGet100(b *testing.B) {
	s := cluster.NewCluster(100, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go setConcurrently(i, s, &wg)
		go getConcurrently(i, s, &wg)
	}
	wg.Wait()
}

func BenchmarkCluster_SetGet1000(b *testing.B) {
	s := cluster.NewCluster(1000, 0)
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go setConcurrently(i, s, &wg)
		go getConcurrently(i, s, &wg)
	}
	wg.Wait()
}
