# sched

an attempt at creating a simple scheduler.

It relies on a priority queue (ex. Binary Heap) to maintain all the jobs needed to be scheduled.

### Features

- [ ] state persistence
- [x] cancel pending job

### In Action

```go
func main() {
	s := sched.New(sched.Config{
		N: 10,
		Q: queue.BinaryHeap(),
		Sink: func(id int, job sched.Job) {
      fmt.Printf("by worker: %d | message: %s\n", id, job.Msg)
		},
	})

  now := time.Now()

	// queue messages
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("msg %d", i)
		s.Queue([]byte(msg), now.Add(3*time.Second))
	}

	jobOneId := s.Queue([]byte("missing!"), now.Add(2*time.Second))

	s.Queue([]byte("You rock!"), now.Add(2*time.Second))

	// cancel job
	s.CancelJob(jobOneId)
}
```
Result:
```txt
by worker: 9 | message: You rock!
by worker: 4 | message: msg 3
by worker: 5 | message: msg 0
by worker: 6 | message: msg 1
by worker: 7 | message: msg 2
by worker: 8 | message: msg 4
```

