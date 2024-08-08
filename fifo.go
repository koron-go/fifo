/*
Package fifo provides simple FIFO (queue) implementation.
*/
package fifo

type entry[T any] struct {
	value T
	prev  *entry[T]
}

// FIFO is a simple implementation of FIFO queue
type FIFO[T any] struct {
	head *entry[T]
	tail *entry[T]
	len  int
}

// Len gets a length of FIFO queue.
func (fifo FIFO[T]) Len() int {
	return fifo.len
}

// Insert inserts a value into FIFO queue.
func (fifo *FIFO[T]) Insert(v T) {
	item := &entry[T]{value: v}
	if fifo.head != nil {
		fifo.head.prev = item
	} else {
		fifo.tail = item
	}
	fifo.head = item
	fifo.len++
}

// Evict evicts the last value from FIFO queue.
func (fifo *FIFO[T]) Evict() (T, bool) {
	if fifo.tail == nil {
		var zero T
		return zero, false
	}
	v := fifo.tail.value
	fifo.tail = fifo.tail.prev
	if fifo.tail == nil {
		fifo.head = nil
	}
	fifo.len--
	return v, true
}
