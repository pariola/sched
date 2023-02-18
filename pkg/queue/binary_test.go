package queue

import (
	"testing"
	"time"
)

// newBH returns basic instance of the bh struct
func newBH() *bh {
	return &bh{
		hChanged: make(chan struct{}, 1),
	}
}

func TestBHEnqueue(t *testing.T) {
	q := newBH()

	now := time.Now()

	jobs := []*Job{
		NewJob(3, now.Add(3*time.Second)),
		NewJob(4, now.Add(4*time.Second)),
		NewJob(31, now.Add(31*time.Second)),
		NewJob(5, now.Add(5*time.Second)),
		NewJob(2, now.Add(2*time.Second)),
	}

	// 2, 3, 31, 5, 4
	expected := []*Job{
		jobs[4], jobs[0], jobs[2], jobs[3], jobs[1],
	}

	// enqueue
	for _, j := range jobs {
		q.Enqueue(j)
	}

	// check size
	if len(q.jobs) != len(jobs) {
		t.Fatalf("queue size/length invalid, expected %d got %d", len(jobs), len(q.jobs))
	}

	if q.Size() != len(jobs) {
		t.Fatalf("queue size/length invalid, expected %d got %d", len(jobs), q.Size())
	}

	// check order
	for i := 0; i < len(jobs); i++ {
		if q.jobs[i] != expected[i] {
			t.Fatalf("order of jobs invalid, index %d - expected id %d got %d", i, expected[i].Id, q.jobs[i].Id)
		}
	}
}

func TestBHDequeue(t *testing.T) {
	q := newBH()

	now := time.Now()

	jobs := []*Job{
		NewJob(3, now.Add(3*time.Second)),
		NewJob(4, now.Add(4*time.Second)),
		NewJob(31, now.Add(31*time.Second)),
		NewJob(5, now.Add(5*time.Second)),
		NewJob(2, now.Add(2*time.Second)),
	}

	// 2, 3, 4, 5, 31
	expectedOrder := []*Job{
		jobs[4], jobs[0], jobs[1], jobs[3], jobs[2],
	}

	// enqueue
	for _, j := range jobs {
		q.Enqueue(j)
	}

	// check order
	for i := 0; i < len(jobs); i++ {

		head := q.Dequeue()

		if head == nil {
			t.Fatalf("job should not be nil")
		}

		if head != expectedOrder[i] {
			t.Fatalf("order of jobs invalid, expected id %d got %d", expectedOrder[i].Id, head.Id)
		}
	}

	if q.Dequeue() != nil {
		t.Fatalf("no job enqueued result should be nil")
	}
}

func TestSatisfyHeapInvariant(t *testing.T) {
	now := time.Now()

	var child, parent *Job

	// fail
	child = NewJob(1, now)
	parent = NewJob(1, now.Add(time.Second))

	if satisfyMinInvariant(parent, child) {
		t.Fatal("min heap invariant should not satisfied")
	}

	// success
	child = NewJob(1, now.Add(time.Second))
	parent = NewJob(1, now)

	if !satisfyMinInvariant(parent, child) {
		t.Fatal("min heap invariant should be satisfied")
	}
}

func TestBHOnHeadChange(t *testing.T) {
	q := newBH()

	enq := make(chan struct{})

	go func() {
		now := time.Now()

		q.Enqueue(NewJob(0, now.Add(7*time.Second)))
		q.Enqueue(NewJob(1, now.Add(3*time.Second)))

		enq <- struct{}{}
	}()

	<-enq

	select {
	case <-q.hChanged:
	default:
		t.Fatal("head change not triggered")
	}
}
