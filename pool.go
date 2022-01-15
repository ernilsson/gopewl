package gopewl

import (
	"fmt"
)

// Job is a simple function that takes no arguments and returns no data. This is the basic operation that can be
// scheduled into the worker pool.
type Job func()

// Pool manages the Job channel and the workers. Currently, there is not much meaning in keeping references to the
// workers, but it is kept here as a way to easily expand the pool with additional functionality requiring knowledge
// of the running workers.
type Pool struct {
	workers []worker
	queue chan Job
}

// Schedule sends the provided Job through to the workers.
func (p *Pool) Schedule(job Job) {
	p.queue <- job
}

// Close simply closes the Job channel, which in turn kills all the worker threads.
func (p *Pool) Close() {
	close(p.queue)
}

// NewPool returns a Pool with the amount of workers specified in `poolSize` and a queue with the buffer size
// of `queueSize`
func NewPool(poolSize int, queueSize int) (*Pool, error) {
	if poolSize <= 0 {
		return nil, fmt.Errorf("invalid pool size '%d'", poolSize)
	}
	if queueSize < 0 {
		return nil, fmt.Errorf("invalid queue size '%d'", queueSize)
	}
	pool := Pool{
		workers: make([]worker, poolSize),
		queue: make(chan Job, queueSize),
	}
	for i := range pool.workers {
		w := newWorker(pool.queue)
		pool.workers[i] = w
		go w.run()
	}
	return &pool, nil
}