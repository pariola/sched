package queue

import "time"

// Item
type Item struct {
	Id   int64
	When time.Time
}

// Queue
type Queue interface {
	Add(*Item)
	Pop() *Item
	Peek() *Item
	Size() int
	Remove(int64)
	OnHeadChange() <-chan struct{}
}

// NewItem
func NewItem(id int64, when time.Time) *Item {
	return &Item{Id: id, When: when}
}
