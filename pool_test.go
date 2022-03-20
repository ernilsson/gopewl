package gopewl

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNewPool_returnsErrorWithNegativePoolSize(t *testing.T) {
	poolSize := -1
	_, err := NewPool(PoolOpts{
		poolSize: poolSize,
	})
	assert.NotNilf(t, err, "Should return error if pool size < 0")
}

func TestNewPool_returnsErrorWithZeroPoolSize(t *testing.T) {
	poolSize := 0
	_, err := NewPool(PoolOpts{
		poolSize: poolSize,
	})
	assert.NotNilf(t, err, "Should return error if pool size = 0")
}

func TestNewPool_returnsErrorWithNegativePoolCapacity(t *testing.T) {
	poolCapacity := -1
	_, err := NewPool(PoolOpts{
		poolSize:     1,
		poolCapacity: poolCapacity,
	})
	assert.NotNilf(t, err, "Should return error if pool capacity < 0")
}

func TestNewPool_returnsErrorWithNonZeroPoolCapacityLargerThanPoolSize(t *testing.T) {
	poolCapacity := 1
	_, err := NewPool(PoolOpts{
		poolSize:     2,
		poolCapacity: poolCapacity,
	})
	assert.NotNilf(t, err, "Should return error if pool capacity > 0 but smaller than pool size")
}

func TestNewPool_returnsNoErrorWithZeroPoolCapacity(t *testing.T) {
	poolCapacity := 0
	_, err := NewPool(PoolOpts{
		poolSize:     1,
		poolCapacity: poolCapacity,
	})
	assert.Nilf(t, err, "Should not return error if pool capacity = 0")
}

func TestNewPool_returnsNoErrorWithPositivePoolSize(t *testing.T) {
	poolSize := 1
	_, err := NewPool(PoolOpts{
		poolSize: poolSize,
	})
	assert.Nilf(t, err, "Should not return error if pool size > 0")
}

func TestNewPool_spawnsCorrectNumberOfGoroutines(t *testing.T) {
	initialNumberOfGoRoutines := runtime.NumGoroutine()
	poolSize := 2
	NewPool(PoolOpts{
		poolSize: poolSize,
	})
	numberOfGoroutines := runtime.NumGoroutine() - initialNumberOfGoRoutines
	assert.Equal(t, poolSize, numberOfGoroutines, "Should spawn one goroutine per worker")
}

func TestNewPool_returnsErrorWithNegativeQueueSize(t *testing.T) {
	queueSize := -1
	_, err := NewPool(PoolOpts{
		poolSize:  1,
		queueSize: queueSize,
	})
	assert.NotNilf(t, err, "Should return error if queue size < 0")
}

func TestNewPool_createsCorrectlyBufferedChannel(t *testing.T) {
	queueSize := 2
	p, _ := NewPool(PoolOpts{
		poolSize:  1,
		queueSize: queueSize,
	})
	channelBuffSize := cap(p.queue)
	assert.Equal(t, queueSize, channelBuffSize, "Should spawn one goroutine per worker")
}

func TestNewPool_wiresWorkersWithChannel(t *testing.T) {
	p, _ := NewPool(PoolOpts{
		poolSize: 2,
	})
	for _, worker := range p.workers {
		assert.Equal(t, p.queue, worker.queue, "Worker should be wired to pool queue upon creation")
	}
}

func TestPool_Schedule_jobIsEventuallyCompletedByWorker(t *testing.T) {
	p, _ := NewPool(PoolOpts{
		poolSize: 2,
	})
	jobIsDone := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	p.Schedule(func() {
		defer wg.Done()
		jobIsDone = true
	})
	wg.Wait()
	assert.Truef(t, jobIsDone, "After scheduling a worker should eventually complete the job")
}

func TestPool_Schedule_createsNewWorkerIfAllWorkersAreOccupiedAndCapacityIsNotReached(t *testing.T) {
	poolSize := 2
	poolCapacity := 4
	p, _ := NewPool(PoolOpts{
		poolSize:     poolSize,
		poolCapacity: poolCapacity,
	})
	for _, worker := range p.workers {
		worker.waiting = false
	}
	p.Schedule(func() {})
	assert.Equal(
		t,
		3,
		len(p.workers),
		"Should create new worker when all others are occupied and capacity is not reached",
	)
}

func TestPool_Schedule_doesNotCreateWorkerIfAllWorkersAreOccupiedAndCapacityIsReached(t *testing.T) {
	poolSize := 2
	poolCapacity := 2
	p, _ := NewPool(PoolOpts{
		poolSize:     poolSize,
		poolCapacity: poolCapacity,
	})
	for _, worker := range p.workers {
		worker.waiting = false
	}
	p.Schedule(func() {})
	assert.Equal(
		t,
		2,
		len(p.workers),
		"Should create new worker when all others are occupied and capacity is not reached",
	)
}

func TestPool_ScheduleWithDelay_waitsBeforeSchedulingJob(t *testing.T) {
	poolSize := 2
	poolCapacity := 2
	p, _ := NewPool(PoolOpts{
		poolSize:     poolSize,
		poolCapacity: poolCapacity,
	})
	beforeScheduleTime := time.Now()
	var afterScheduleTime time.Time
	wg := sync.WaitGroup{}
	wg.Add(1)
	p.ScheduleWithDelay(func() {
		afterScheduleTime = time.Now()
		wg.Done()
	}, time.Second)
	wg.Wait()

	assert.NotNilf(t, afterScheduleTime, "Job did not execute correctly")
	assert.Truef(t, afterScheduleTime.Sub(beforeScheduleTime).Seconds() >= 1, "Job did not run one second after scheduling.")
	assert.Truef(t, afterScheduleTime.Sub(beforeScheduleTime).Seconds() < 2, "Job took more than one second to run after scheduling.")
}

func TestPool_Close_closesJobQueue(t *testing.T) {
	p, _ := NewPool(PoolOpts{
		poolSize: 2,
	})
	p.Close()
	var queueIsClosed bool
	select {
	case <-p.queue:
		queueIsClosed = true
	default:
		queueIsClosed = false
	}
	assert.True(t, queueIsClosed, "Queue should be closed when Pool is")
}
