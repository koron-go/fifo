package fifo

import "iter"

// Seq returns iter.Seq[T] to enumerate items in FIFO queue.
func (fifo *FIFO[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for p := fifo.tail; p != nil; p = p.prev {
			if !yield(p.value) {
				break
			}
		}
	}
}
