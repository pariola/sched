package queue

import (
	"sched/types"
	"testing"
	"time"
)

func TestBHEnqueue(t *testing.T) {

	var q bh

	now := time.Now()

	jobs := []types.Job{
		{Id: 3, When: now.Add(3 * time.Second)},
		{Id: 4, When: now.Add(4 * time.Second)},
		{Id: 31, When: now.Add(31 * time.Second)},
		{Id: 5, When: now.Add(5 * time.Second)},
		{Id: 2, When: now.Add(2 * time.Second)},
	}

	// 2, 3, 31, 5, 4
	expected := []types.Job{
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

	var q bh

	now := time.Now()

	jobs := []types.Job{
		{Id: 3, When: now.Add(3 * time.Second)},
		{Id: 4, When: now.Add(4 * time.Second)},
		{Id: 31, When: now.Add(31 * time.Second)},
		{Id: 5, When: now.Add(5 * time.Second)},
		{Id: 2, When: now.Add(2 * time.Second)},
	}

	// 2, 3, 4, 5, 31
	expectedOrder := []types.Job{
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

		if *head != expectedOrder[i] {
			t.Fatalf("order of jobs invalid, expected id %d got %d", expectedOrder[i].Id, head.Id)
		}
	}

	if q.Dequeue() != nil {
		t.Fatalf("no job enqueued result should be nil")
	}
}

func TestSatisfyHeapInvariant(t *testing.T) {

	now := time.Now()

	var child, parent types.Job

	// fail
	child = types.Job{Id: 1, When: now}
	parent = types.Job{Id: 1, When: now.Add(time.Second)}

	if satisfyMinInvariant(parent, child) {
		t.Fatal("min heap invariant should not satisfied")
	}

	// success
	child = types.Job{Id: 1, When: now.Add(time.Second)}
	parent = types.Job{Id: 1, When: now}

	if !satisfyMinInvariant(parent, child) {
		t.Fatal("min heap invariant should be satisfied")
	}
}
