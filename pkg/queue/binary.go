package queue

import (
	"sched/pkg/types"
	"sync"
)

// bh - binary heap
type bh struct {
	m sync.Mutex

	// jobs
	jobs []*types.Job

	// hChanged signals on head change
	hChanged chan struct{}
}

// New returns an instance of the binary heap
func NewBH() Queue {
	return &bh{
		hChanged: make(chan struct{}, 1),
	}
}

// satisfyMinInvariant compares parent and child to verify min heap invariant
func satisfyMinInvariant(parent, child *types.Job) bool {
	return parent.When.Before(child.When)
}

// Size returns the amount of jobs enqueued
func (h *bh) Size() int {
	return len(h.jobs)
}

// Enqueue enqueues a job
func (h *bh) Enqueue(job *types.Job) {

	h.m.Lock()
	defer h.m.Unlock()

	// queue job
	h.jobs = append(h.jobs, job)

	// head
	head := h.jobs[0]

	// bubble-up
	pos := len(h.jobs) - 1

	for pos > 0 {

		// parent position
		ppos := (pos - 1) / 2

		// compare current and parent
		if satisfyMinInvariant(h.jobs[ppos], h.jobs[pos]) {
			break
		}

		// swap current <=> parent
		h.jobs[pos], h.jobs[ppos] = h.jobs[ppos], h.jobs[pos]

		// update position
		pos = ppos
	}

	// signal on head change
	if head != h.jobs[0] {
		select {
		case h.hChanged <- struct{}{}:
		case <-h.hChanged:
		}
	}
}

// Dequeue returns the next enqueued job
func (h *bh) Dequeue() *types.Job {

	h.m.Lock()
	defer h.m.Unlock()

	// no jobs
	if len(h.jobs) < 1 {
		return nil
	}

	var pos int

	head := h.jobs[pos]

	bottom := len(h.jobs) - 1

	// replace with bottom
	h.jobs[pos] = h.jobs[bottom]

	// remove last job
	h.jobs = h.jobs[:bottom]

	// bubble down
	for {
		l, r := (2*pos)+1, (2*pos)+2

		// out of bounds
		if l >= len(h.jobs) {
			break
		}

		// smallest child node
		child := l
		if r < len(h.jobs) && h.jobs[l].When.After(h.jobs[r].When) {
			child = r
		}

		// verify heap invariant
		if satisfyMinInvariant(h.jobs[pos], h.jobs[child]) {
			break
		}

		// swap current <=> parent
		h.jobs[pos], h.jobs[child] = h.jobs[child], h.jobs[pos]

		// update position
		pos = child
	}

	return head
}

// Peek returns the next enqueued job without removing it from queue
func (h *bh) Peek() *types.Job {

	h.m.Lock()
	defer h.m.Unlock()

	// no jobs
	if len(h.jobs) < 1 {
		return nil
	}

	return h.jobs[0]
}

// OnHeadChange returns a channel to listen to queue head change events
func (h *bh) OnHeadChange() <-chan struct{} {
	return h.hChanged
}
