package gopewl

// worker holds a reference to the Job channel and continuously listens for jobs to complete.
type worker struct {
	queue chan Job
}

// run is run as a goroutine, and its only responsibility is to read jobs from the Job channel and execute
// them.
func (w worker) run() {
	for job := range w.queue {
		job()
	}
}

// newWorker returns a worker with the provided Job channel.
func newWorker(queue chan Job) worker {
	return worker{
		queue: queue,
	}
}
