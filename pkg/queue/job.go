package queue

import "time"

type Job struct {
	Id   int64
	When time.Time
}

func NewJob(id int64, when time.Time) *Job {
	return &Job{Id: id, When: when}
}
