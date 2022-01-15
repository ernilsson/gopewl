package gopewl

import (
	"errors"
)

// Job is a simple function that takes no arguments and returns no data. This is the basic operation that can be
// scheduled into the worker pool.
type Job func()

// Pool manages the Job channel and the workers. It ensures that the number of workers are as expected and spawns
// or kills workers as needed.
type Pool struct {
	capacity int
	workers []worker
	queue chan Job
}

// Schedule sends the provided Job through to the workers.
func (p *Pool) Schedule(job Job) {
	if p.allWorkersAreOccupied() && p.canAddWorker() {
		p.addWorker()
	}
	p.queue <- job
}

func (p *Pool) addWorker() {
	w := newWorker(p.queue)
	p.workers = append(p.workers, w)
	go w.run()
}

func (p Pool) allWorkersAreOccupied() bool {
	for _, worker := range p.workers {
		if worker.waiting {
			return false
		}
	}
	return true
}

func (p Pool) canAddWorker() bool {
	return p.capacity > 0 && p.capacity < len(p.workers)
}

// Close simply closes the Job channel, which in turn kills all the worker threads.
func (p *Pool) Close() {
	close(p.queue)
}

type PoolOptions struct {
	poolSize int
	poolCapacity int
	queueSize int
}

func (po PoolOptions) validate() error {
	if po.poolSize <= 0 {
		return errors.New("pool size must be a positive integer")
	}
	if po.queueSize < 0 {
		return errors.New("queue size must be a positive integer or 0")
	}
	if po.poolCapacity > 0 && po.poolCapacity < po.poolSize {
		return errors.New("pool capacity must be larger than pool size or 0")
	}
	if po.poolCapacity < 0 {
		return errors.New("pool capacity must be a positive integer or 0")
	}
	return nil
 }

// NewPool returns a Pool with the amount of workers specified in `poolSize` and a queue with the buffer size
// of `queueSize`
func NewPool(opts PoolOptions) (*Pool, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	p := &Pool{
		capacity: opts.poolCapacity,
		workers: make([]worker, 0),
		queue: make(chan Job, opts.queueSize),
	}
	for i := 0; i < opts.poolSize; i++ {
		p.addWorker()
	}
	return p, nil
}