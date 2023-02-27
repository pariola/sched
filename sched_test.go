package sched_test

import (
	"testing"
	"time"

	"github.com/pariola/sched"
	"github.com/pariola/sched/internal/pkg/queue"
)

func TestSched(t *testing.T) {
	done := make(chan sched.Job)

	s := sched.New(sched.Config{
		N: 1,
		Q: queue.BinaryHeap(),
		Sink: func(id int, job sched.Job) {
			done <- job
		},
	})

	msg := []byte("Hi there")
	s.Queue(msg, time.Now().Add(2*time.Second))

	select {
	case j := <-done:
		if string(j.Msg) != string(msg) {
			t.Fatalf("expected msg: %+v, got: %+v\n", msg, j.Msg)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("trigger timeout")
	}
}
