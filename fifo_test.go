package fifo_test

import (
	"testing"

	"github.com/koron-go/fifo"
)

func TestFIFO(t *testing.T) {
	var fifo fifo.FIFO[int]
	if d := fifo.Len(); d != 0 {
		t.Fatalf("fifo length is not zero: got=%d", d)
	}

	// Insert
	fifo.Insert(1)
	fifo.Insert(2)
	fifo.Insert(3)
	if d := fifo.Len(); d != 3 {
		t.Fatalf("fifo length is not 3: got=%d", d)
	}
	fifo.Insert(4)
	fifo.Insert(5)
	if d := fifo.Len(); d != 5 {
		t.Fatalf("fifo length is not 5: got=%d", d)
	}

	// Evict
	if v, ok := fifo.Evict(); !ok || v != 1 {
		t.Fatalf("failed to evict want=1: ok=%t got=%d", ok, v)
	}
	if d := fifo.Len(); d != 4 {
		t.Fatalf("fifo length is not 4: got=%d", d)
	}
	if v, ok := fifo.Evict(); !ok || v != 2 {
		t.Fatalf("failed to evict want=2: ok=%t got=%d", ok, v)
	}
	if v, ok := fifo.Evict(); !ok || v != 3 {
		t.Fatalf("failed to evict want=3: ok=%t got=%d", ok, v)
	}
	if v, ok := fifo.Evict(); !ok || v != 4 {
		t.Fatalf("failed to evict want=4: ok=%t got=%d", ok, v)
	}
	if d := fifo.Len(); d != 1 {
		t.Fatalf("fifo length is not 1: got=%d", d)
	}
	if v, ok := fifo.Evict(); !ok || v != 5 {
		t.Fatalf("failed to evict want=5: ok=%t got=%d", ok, v)
	}
	if d := fifo.Len(); d != 0 {
		t.Fatalf("fifo length is not 0: got=%d", d)
	}

	// Evict from empty FIFO queue.
	if v, ok := fifo.Evict(); ok || v != 0 {
		t.Fatalf("unexpectedly a value evicted from empty FIFO queue: ok=%t got=%d", ok, v)
	}

	// Insert and Evict again
	fifo.Insert(6)
	fifo.Insert(7)
	fifo.Insert(8)
	if d := fifo.Len(); d != 3 {
		t.Fatalf("fifo length is not 3: got=%d", d)
	}
	if v, ok := fifo.Evict(); !ok || v != 6 {
		t.Fatalf("failed to evict want=6: ok=%t got=%d", ok, v)
	}
	if v, ok := fifo.Evict(); !ok || v != 7 {
		t.Fatalf("failed to evict want=7: ok=%t got=%d", ok, v)
	}
	if v, ok := fifo.Evict(); !ok || v != 8 {
		t.Fatalf("failed to evict want=8: ok=%t got=%d", ok, v)
	}
	if d := fifo.Len(); d != 0 {
		t.Fatalf("fifo length is not 0: got=%d", d)
	}
}
