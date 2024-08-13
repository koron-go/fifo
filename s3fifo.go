package fifo

import (
	"bytes"
	"fmt"
	"log"
)

type s3item[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

// S3FIFO implements S3-FIFO.
// See also https://s3fifo.com/
type S3FIFO[K comparable, V any] struct {
	small FIFO[s3item[K, V]]
	main  FIFO[s3item[K, V]]
	ghost FIFO[K]

	smallSize int
	mainSize  int
}

func NewS3FIFO[K comparable, V any](smallSize, mainSize int) *S3FIFO[K, V] {
	if smallSize < 1 {
		smallSize = 1
	}
	if mainSize < smallSize {
		mainSize = smallSize
	}
	return &S3FIFO[K, V]{
		smallSize: smallSize,
		mainSize:  mainSize,
	}
}

func (cache *S3FIFO[K, V]) Size() int {
	return cache.smallSize + cache.mainSize
}

func (cache *S3FIFO[K, V]) Len() int {
	return cache.small.Len() + cache.main.Len()
}

func (cache *S3FIFO[K, V]) inSmallOrMain(k K) (*s3item[K, V], bool) {
	equalKey := func(item s3item[K, V]) bool {
		return item.key == k
	}
	if p, ok := cache.small.Find(equalKey); ok {
		return p, true
	}
	return cache.main.Find(equalKey)
}

func (cache *S3FIFO[K, V]) Get(k K) (V, bool) {
	item, ok := cache.inSmallOrMain(k)
	if !ok {
		var zero V
		return zero, false
	}
	if item.freq < 3 {
		item.freq++
	}
	return item.value, true
}

func (cache *S3FIFO[K, V]) Put(k K, v V) {
	item, ok := cache.inSmallOrMain(k)
	if ok {
		// update a cached item.
		item.value = v
		if item.freq < 3 {
			item.freq++
		}
		return
	}
	cache.insertNew(s3item[K, V]{key: k, value: v, freq: 0})
}

func (cache *S3FIFO[K, V]) insertNew(item s3item[K, V]) {
	if cache.removeFromGhost(item.key) {
		log.Printf("removed from ghost, move to main: key=%+v", item.key)
		cache.evictM()
		cache.main.Insert(item)
	} else {
		cache.evictS()
		cache.small.Insert(item)
	}
}

func (cache *S3FIFO[K, V]) removeFromGhost(k K) bool {
	return cache.ghost.RemoveIf(func(v K) bool {
		return v == k
	})
}

func (cache *S3FIFO[K, V]) evictM() {
	for cache.main.Len() >= cache.mainSize {
		item, _ := cache.main.Evict()
		if item.freq > 0 {
			item.freq--
			cache.main.Insert(item)
		} else {
			log.Printf("evicted from main: key=%+v value=%+v", item.key, item.value)
		}
	}
}

func (cache *S3FIFO[K, V]) evictS() {
	for cache.small.Len() >= cache.smallSize {
		item, _ := cache.small.Evict()
		if item.freq > 0 {
			cache.evictM()
			item.freq = 0
			cache.main.Insert(item)
			log.Printf("evicted from small, move to main: item=%+v", item)
		} else {
			for cache.ghost.Len() >= cache.mainSize {
				cache.ghost.Evict()
			}
			cache.ghost.Insert(item.key)
			log.Printf("evicted from small, move to ghost, value dropped: item=%+v", item)
		}
	}
}

func enumerateFIFO[V any](e *entry[V], fn func(v V)) {
	if e == nil {
		return
	}
	enumerateFIFO(e.prev, fn)
	fn(e.value)
}

func (cache *S3FIFO[K, V]) Dump() {
	bb := &bytes.Buffer{}
	fmt.Fprintf(bb, "s[")
	enumerateFIFO(cache.small.tail, func(item s3item[K, V]) {
		fmt.Fprintf(bb, " %+v", item)
	})
	fmt.Fprintf(bb, "]\n")
	fmt.Fprintf(bb, "\t\t    m[")
	enumerateFIFO(cache.main.tail, func(item s3item[K, V]) {
		fmt.Fprintf(bb, " %+v", item)
	})
	fmt.Fprintf(bb, "]\n")
	fmt.Fprintf(bb, "\t\t    g[")
	enumerateFIFO(cache.ghost.tail, func(k K) {
		fmt.Fprintf(bb, " %+v", k)
	})
	fmt.Fprintf(bb, "]\n")
	log.Print(bb.String())
}
