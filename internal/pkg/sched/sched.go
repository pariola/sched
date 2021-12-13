package sched

import (
	"sched/pkg/queue"
	"sched/pkg/types"
	"time"
)

type Sink func(int, *types.Job)

type backend struct {

	// q the pending jobs queue
	q queue.Queue

	// jobs the processing queue
	jobs chan *types.Job

	// stopDrainC the stop signal channel for drain-ing jobs from q
	stopDrainC chan struct{}

	// workerStopC contains stop signal channel for all workers
	workerStopC []chan struct{}
}

func New() *backend {

	b := &backend{
		q:          queue.NewBH(),
		stopDrainC: make(chan struct{}),
		jobs:       make(chan *types.Job, 1000),
	}

	// drain
	go func() {
		for {
			head := b.q.Peek()

			if head == nil {
				continue
			}

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
			case <-b.q.OnHeadChange():
				timer.Stop()
				continue
			}

			// send to processing queue
			b.jobs <- head

			time.Sleep(10 * time.Millisecond)
		}
	}()

	return b
}

// worker
func worker(id int, jobs chan *types.Job, stopC chan struct{}, sink Sink) {
	for {
		select {
		case <-stopC:
			return
		case j := <-jobs:
			sink(id, j)
		}
	}
}

// Sink creates neccessary workers
func (b *backend) Sink(n int, sink Sink) {
	b.workerStopC = make([]chan struct{}, n)

	for i := 0; i < n; i++ {
		b.workerStopC[i] = make(chan struct{})
		go worker(i, b.jobs, b.workerStopC[i], sink)
	}
}

// Stop signals the scheduler to stop draining jobs
func (b *backend) Stop() {

	// stop drain
	b.stopDrainC <- struct{}{}
	close(b.stopDrainC)

	// stop workers
	for _, c := range b.workerStopC {
		c <- struct{}{}
		close(c)
	}

	// close processing queue
	close(b.jobs)
}

// CreateJob enqueues a new pending job
func (b *backend) CreateJob(job *types.Job) {
	b.q.Enqueue(job)
}

// UpdateJob
func (b *backend) UpdateJob(id int64) {
}

// CancelJob
func (b *backend) CancelJob(id int64) {
}
