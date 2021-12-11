package types

import (
	"time"
)

type Job struct {
	Id int64 `json:"id"`

	When time.Time `json:"when"`
}
