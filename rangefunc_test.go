package fifo_test

import (
	"testing"

	"github.com/koron-go/fifo"
)

func TestFIFO_Seq(t *testing.T) {
	var fifo fifo.FIFO[int]
	fifo.Insert(1)
	fifo.Insert(2)
	fifo.Insert(3)
	fifo.Insert(4)
	fifo.Insert(5)
	wants := []int{1, 2, 3, 4, 5}
	for got := range fifo.Seq() {
		if len(wants) == 0 {
			t.Fatalf("wants exhausted: got=%d", got)
		}
		want := wants[0]
		wants = wants[1:]
		if got != want {
			t.Fatalf("unexpected value: want=%d got=%d", want, got)
		}
	}
	if len(wants) != 0 {
		t.Fatalf("wants remained: %+v", wants)
	}
	if d := fifo.Len(); d != 5 {
		t.Fatalf("fifo length is not 5: got=%d", d)
	}
}

func TestFIFO_SeqBreak(t *testing.T) {
	var fifo fifo.FIFO[int]
	fifo.Insert(1)
	fifo.Insert(2)
	fifo.Insert(3)
	fifo.Insert(4)
	fifo.Insert(5)

	wants := []int{1, 2, 3}
	for got := range fifo.Seq() {
		if len(wants) == 0 {
			break
		}
		want := wants[0]
		wants = wants[1:]
		if got != want {
			t.Fatalf("unexpected value: want=%d got=%d", want, got)
		}
	}
	if len(wants) != 0 {
		t.Fatalf("wants remained: %+v", wants)
	}
	if d := fifo.Len(); d != 5 {
		t.Fatalf("fifo length is not 5: got=%d", d)
	}
}
