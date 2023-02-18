package queue

// Queue
type Queue interface {
	Size() int
	Enqueue(*Job)
	Dequeue() *Job
	Peek() *Job
	OnHeadChange() <-chan struct{}
}
