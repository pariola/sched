package queue

import (
	"sched/pkg/types"
)

// Queue
type Queue interface {
	Size() int
	Enqueue(*types.Job)
	Dequeue() *types.Job
	Peek() *types.Job
	OnHeadChange() <-chan struct{}
}
