package sched

import (
	"time"

	"github.com/pariola/sched/internal/pkg/queue"
)

// Job
type Job struct {
	Id  int64
	Msg []byte
}

// Sink actor
type Sink func(id int, job Job)

type Config struct {
	N    int
	Q    queue.Queue
	Sink Sink
}

// backend
type backend struct {
	// q the pending jobs queue
	q queue.Queue

	// items represents the processing queue
	items chan *queue.Item

	// storage - a map of job id to job message
	storage map[int64][]byte

	// stopC the stop signal channel for drain-ing jobs from q
	stopC chan struct{}

	// workerStopC contains stop signal channel for all workers
	workerStopC []chan struct{}
}

// New returns an instance of backend
func New(c Config) *backend {
	if c.N < 1 {
		c.N = 1
	}

	b := &backend{
		q:           c.Q,
		stopC:       make(chan struct{}),
		storage:     make(map[int64][]byte),
		items:       make(chan *queue.Item, 1000),
		workerStopC: make([]chan struct{}, c.N),
	}

	// start drain process
	go b.processor()

	for i := 0; i < c.N; i++ {
		b.workerStopC[i] = make(chan struct{})
		go b.worker(i, b.workerStopC[i], c.Sink)
	}

	return b
}

// worker
func (b *backend) worker(id int, stopC chan struct{}, sink Sink) {
	for {
		select {
		case <-stopC:
			return
		case i := <-b.items:
			sink(id, Job{Id: i.Id, Msg: b.storage[i.Id]})
		}
	}
}

// processor
func (b *backend) processor() {
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
			head = b.q.Pop()

		case <-b.stopC:
			timer.Stop()
			return

		// break out as soon as head changes
		case <-b.q.OnHeadChange():
			timer.Stop()
			continue
		}

		b.items <- head // send to processing queue

		time.Sleep(10 * time.Millisecond)
	}
}

// Stop signals the scheduler and workers to stop
func (b *backend) Stop() {
	b.stopC <- struct{}{} // stop drain
	close(b.stopC)

	// stop workers
	for _, c := range b.workerStopC {
		c <- struct{}{}
		close(c)
	}

	close(b.items) // close processing queue
}

// Queue creates a new job
func (b *backend) Queue(msg []byte, when time.Time) int64 {
	id := time.Now().Unix() + int64(len(b.storage))

	b.storage[id] = msg

	b.q.Add(queue.NewItem(id, when))

	return id
}

// CancelJob cancels a pending job
func (b *backend) CancelJob(id int64) {
	b.q.Remove(id)
	delete(b.storage, id)
}
