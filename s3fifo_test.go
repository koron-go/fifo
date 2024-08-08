package fifo_test

import (
	"testing"

	"github.com/koron-go/fifo"
)

func TestS3FIFO(t *testing.T) {
	cache := fifo.NewS3FIFO[int, string](2, 10)
	cache.Dump()
	cache.Put(1, "a1")
	cache.Dump()
	cache.Put(2, "b1")
	cache.Dump()
	cache.Put(1, "a2")
	cache.Dump()
	cache.Put(3, "c1")
	cache.Dump()
	cache.Put(4, "d1")
	cache.Dump()
	cache.Put(5, "e1")
	cache.Dump()
	cache.Put(1, "a3")
	cache.Dump()
	cache.Put(3, "c2")
	cache.Dump()
}
