package queue

import (
	"testing"
	"time"
)

// newBH returns basic instance of the bh struct
func newBH() *bh {
	return &bh{
		hC: make(chan struct{}, 1),
	}
}

func TestBinaryHeapAdd(t *testing.T) {
	q := newBH()

	now := time.Now()

	items := [...]*Item{
		NewItem(3, now.Add(3*time.Second)),
		NewItem(4, now.Add(4*time.Second)),
		NewItem(31, now.Add(31*time.Second)),
		NewItem(5, now.Add(5*time.Second)),
		NewItem(2, now.Add(2*time.Second)),
	}

	// 2, 3, 31, 5, 4
	expected := [...]*Item{
		items[4], items[0], items[2], items[3], items[1],
	}

	// enqueue
	for _, j := range items {
		q.Add(j)
	}

	// check size
	if len(q.items) != len(items) {
		t.Fatalf("queue size/length invalid, expected %d got %d", len(items), len(q.items))
	}

	if s := q.Size(); s != len(items) {
		t.Fatalf("queue size/length invalid, expected %d got %d", len(items), s)
	}

	// check order
	for i := 0; i < len(items); i++ {
		if q.items[i] != expected[i] {
			t.Fatalf("order of items invalid, index %d - expected id %d got %d", i, expected[i].Id, q.items[i].Id)
		}
	}
}

func TestBinaryHeapPop(t *testing.T) {
	q := newBH()

	now := time.Now()

	items := [...]*Item{
		NewItem(3, now.Add(3*time.Second)),
		NewItem(4, now.Add(4*time.Second)),
		NewItem(31, now.Add(31*time.Second)),
		NewItem(5, now.Add(5*time.Second)),
		NewItem(2, now.Add(2*time.Second)),
	}

	// 2, 3, 4, 5, 31
	expectedOrder := [...]*Item{
		items[4], items[0], items[1], items[3], items[2],
	}

	// enqueue
	for _, j := range items {
		q.Add(j)
	}

	// check order
	for i := 0; i < len(items); i++ {
		head := q.Pop()

		if head == nil {
			t.Fatalf("job should not be nil")
		}

		if head != expectedOrder[i] {
			t.Fatalf("order of items invalid, expected id %d got %d", expectedOrder[i].Id, head.Id)
		}
	}

	if q.Pop() != nil {
		t.Fatalf("no job enqueued result should be nil")
	}
}

func TestBinaryHeapRemove(t *testing.T) {
	q := newBH()

	now := time.Now()

	items := [...]*Item{
		NewItem(3, now.Add(3*time.Second)),
		NewItem(4, now.Add(4*time.Second)),
		NewItem(31, now.Add(31*time.Second)),
		NewItem(5, now.Add(5*time.Second)),
		NewItem(2, now.Add(2*time.Second)),
	}

	// 3, 4, 31
	expectedOrder := [...]*Item{
		items[0], items[1], items[2],
	}

	// enqueue
	for _, j := range items {
		q.Add(j)
	}

	// remove items 2 & 5
	q.Remove(2)
	q.Remove(5)

	if a, b := q.Size(), len(expectedOrder); a != b {
		t.Fatalf("queue size/length invalid, expected %d got %d", b, a)
	}

	// check order
	for i := 0; i < len(expectedOrder); i++ {
		head := q.Pop()

		if head == nil {
			t.Fatalf("job should not be nil")
		}

		if head != expectedOrder[i] {
			t.Fatalf("order of items invalid, expected id %d got %d", expectedOrder[i].Id, head.Id)
		}
	}

	if q.Pop() != nil {
		t.Fatalf("no job enqueued result should be nil")
	}
}

func TestSatisfyHeapInvariant(t *testing.T) {
	now := time.Now()

	var child, parent *Item

	// fail
	child = NewItem(1, now)
	parent = NewItem(1, now.Add(time.Second))

	if satisfyMinInvariant(parent, child) {
		t.Fatal("min heap invariant should not satisfied")
	}

	// success
	child = NewItem(1, now.Add(time.Second))
	parent = NewItem(1, now)

	if !satisfyMinInvariant(parent, child) {
		t.Fatal("min heap invariant should be satisfied")
	}
}

func TestBinaryHeapOnHeadChange(t *testing.T) {
	q := newBH()

	enq := make(chan struct{})

	go func() {
		now := time.Now()

		q.Add(NewItem(0, now.Add(7*time.Second)))
		q.Add(NewItem(1, now.Add(3*time.Second)))

		close(enq)
	}()

	<-enq

	select {
	case <-q.hC:
	default:
		t.Fatal("head change not triggered")
	}
}
