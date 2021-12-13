package types

import (
	"time"
)

type Job struct {
	Id int64 `json:"id"`

	When time.Time `json:"when"`
}

func NewJob(id int64, when time.Time) *Job {
	return &Job{Id: id, When: when}
}
