package gopewl

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
)

func TestNewPool_returnsErrorWithNegativePoolSize(t *testing.T) {
	poolSize := -1
	_, err := NewPool(poolSize, 0)
	assert.NotNilf(t, err, "Should return error if pool size < 0")
}

func TestNewPool_returnsErrorWithZeroPoolSize(t *testing.T) {
	poolSize := 0
	_, err := NewPool(poolSize, 0)
	assert.NotNilf(t, err, "Should return error if pool size = 0")
}

func TestNewPool_returnsNoErrorWithPositivePoolSize(t *testing.T) {
	poolSize := 1
	_, err := NewPool(poolSize, 0)
	assert.Nilf(t, err, "Should not return error if pool size > 0")
}

func TestNewPool_spawnsCorrectNumberOfGoroutines(t *testing.T) {
	initialNumberOfGoRoutines := runtime.NumGoroutine()
	expectedPoolSize := 2
	NewPool(expectedPoolSize, 0)
	numberOfGoroutines := runtime.NumGoroutine() - initialNumberOfGoRoutines
	assert.Equal(t, expectedPoolSize, numberOfGoroutines, "Should spawn one goroutine per worker")
}

func TestNewPool_returnsErrorWithNegativeQueueSize(t *testing.T) {
	queueSize := -1
	_, err := NewPool(1, queueSize)
	assert.NotNilf(t, err, "Should return error if queue size < 0")
}

func TestNewPool_createsCorrectlyBufferedChannel(t *testing.T) {
	expectedChannelBuffSize := 2
	p, _ := NewPool(1, expectedChannelBuffSize)
	channelBuffSize := cap(p.queue)
	assert.Equal(t, expectedChannelBuffSize, channelBuffSize, "Should spawn one goroutine per worker")
}

func TestNewPool_wiresWorkersWithChannel(t *testing.T) {
	p, _ := NewPool(2, 0)
	for _, worker := range p.workers {
		assert.Equal(t, p.queue, worker.queue, "Worker should be wired to pool queue upon creation")
	}
}

func TestPool_Schedule_jobIsEventuallyCompletedByWorker(t *testing.T) {
	p, _ := NewPool(2, 0)
	jobIsDone := false
	wg := sync.WaitGroup{}
	wg.Add(1)
	p.Schedule(func () {
		defer wg.Done()
		jobIsDone = true
	})
	wg.Wait()
	assert.Truef(t, jobIsDone, "After scheduling a worker should eventually complete the job")
}

func TestPool_Close_closesJobQueue(t *testing.T) {
	p, _ := NewPool(2, 0)
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