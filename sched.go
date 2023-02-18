package sched

import (
	"time"

	"github.com/pariola/sched/pkg/queue"
)

// Sink actor
type Sink func(id int, job *queue.Job)

// backend
type backend struct {
	// q the pending jobs queue
	q queue.Queue

	// jobs the processing queue
	jobs chan *queue.Job

	// stopDrainC the stop signal channel for drain-ing jobs from q
	stopDrainC chan struct{}

	// workerStopC contains stop signal channel for all workers
	workerStopC []chan struct{}
}

// New returns an instance of backend
func New(q queue.Queue) *backend {
	b := &backend{
		q:          q,
		stopDrainC: make(chan struct{}),
		jobs:       make(chan *queue.Job, 1000),
	}

	// drain
	go func() {
		for {
			head := b.q.Peek()

			if head == nil {
				continue
			}

			// block until the head's scheduled time
			timer := time.NewTimer(
				time.Until(head.When),
			)

			select {
			case <-timer.C:
				timer.Stop()
				head = b.q.Dequeue()

			case <-b.stopDrainC:
				timer.Stop()
				return

				// break out as soon as head changes
			case <-b.q.OnHeadChange():
				timer.Stop()
				continue
			}

			b.jobs <- head // send to processing queue

			time.Sleep(10 * time.Millisecond)
		}
	}()

	return b
}

// worker
func worker(id int, jobs chan *queue.Job, stopC chan struct{}, sink Sink) {
	for {
		select {
		case <-stopC:
			return
		case j := <-jobs:
			sink(id, j)
		}
	}
}

// Sink creates n workers with sink actor
func (b *backend) Sink(n int, sink Sink) {
	b.workerStopC = make([]chan struct{}, n)

	for i := 0; i < n; i++ {
		b.workerStopC[i] = make(chan struct{})
		go worker(i, b.jobs, b.workerStopC[i], sink)
	}
}

// Stop signals the scheduler and workers to stop
func (b *backend) Stop() {
	b.stopDrainC <- struct{}{} // stop drain
	close(b.stopDrainC)

	// stop workers
	for _, c := range b.workerStopC {
		c <- struct{}{}
		close(c)
	}

	close(b.jobs) // close processing queue
}

// CreateJob creates a new job
func (b *backend) CreateJob(job *queue.Job) {
	b.q.Enqueue(job)
}

// CancelJob cancels a pending job
func (b *backend) CancelJob(id int64) {
	panic("not implemented")
}
