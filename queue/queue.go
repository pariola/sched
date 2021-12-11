package queue

import "sched/types"

// Queue
type Queue interface {
	Size() int
	Enqueue(types.Job)
	Dequeue() *types.Job
}
