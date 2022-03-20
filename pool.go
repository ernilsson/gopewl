package gopewl

import "time"

// Job is a simple function that takes no arguments and returns no data. This is the basic operation that can be
// scheduled into the worker pool.
type Job func()

// Pool manages the Job channel and the workers. It ensures that the number of workers are as expected and spawns
// or kills workers as needed.
type Pool struct {
	capacity int
	workers  []worker
	queue    chan Job
}

// Schedule sends the provided Job through to the workers.
func (p *Pool) Schedule(job Job) {
	if p.allWorkersAreOccupied() && p.canAddWorker() {
		p.addWorker()
	}
	p.queue <- job
}

// ScheduleWithDelay returns immediately and schedules the job to be run after the given time.Duration delay.
// Beware that this does spawn an extra goroutine with the responsibility of scheduling the job after the delay.
func (p *Pool) ScheduleWithDelay(job Job, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		p.Schedule(job)
	}()
}

// addWorker spawns a new worker and wires it up to the Pool
func (p *Pool) addWorker() {
	w := newWorker(p.queue)
	p.workers = append(p.workers, w)
	go w.run()
}

// allWorkersAreOccupied returns true if all current workers are busy processing a job at the time of calling
func (p Pool) allWorkersAreOccupied() bool {
	for _, worker := range p.workers {
		if worker.waiting {
			return false
		}
	}
	return true
}

// canAddWorker returns true if the Pool has an explicit capacity that has not been reached
func (p Pool) canAddWorker() bool {
	return p.capacity > 0 && len(p.workers) < p.capacity
}

// Close simply closes the Job channel, which in turn kills all the worker threads.
func (p *Pool) Close() {
	close(p.queue)
}

// PoolOpts contains fields used to construct a Pool type.
type PoolOpts struct {
	// poolSize refers to the minimum amount of workers that will always be active until the Pool is closed
	poolSize int
	// poolCapacity refers to the maximum amount of workers that the Pool can utilise
	poolCapacity int
	// queueSize refers to the maximum unprocessed jobs that cab be present before the scheduling of a new Job will
	// become blocking
	queueSize int
}

// validate returns an error if any of the PoolOpts fields contains illegal values
func (po PoolOpts) validate() error {
	if po.poolSize <= 0 {
		return ErrNonPositivePoolSize
	}
	if po.queueSize < 0 {
		return ErrNegativeQueueSize
	}
	if po.poolCapacity > 0 && po.poolCapacity < po.poolSize {
		return ErrIllegalPoolCapacity
	}
	if po.poolCapacity < 0 {
		return ErrNegativePoolCapacity
	}
	return nil
}

// NewPool returns a Pool with the amount of workers specified in `poolSize` and a queue with the buffer size
// of `queueSize`
func NewPool(opts PoolOpts) (*Pool, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	p := &Pool{
		capacity: opts.poolCapacity,
		workers:  make([]worker, 0),
		queue:    make(chan Job, opts.queueSize),
	}
	for i := 0; i < opts.poolSize; i++ {
		p.addWorker()
	}
	return p, nil
}
