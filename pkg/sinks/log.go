package sinks

import (
	"log"
	"sched/pkg/types"
)

// LogSink logs the worker running the current job
func LogSink(wid int, j *types.Job) {
	log.Println("Worker", wid, "Running Job", j.Id)
}
