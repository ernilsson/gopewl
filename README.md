# gopewl 🌊
A simple package allowing quick capping of goroutines for rudimentary asynchronous tasks.

## Example
The API is trivial to use. Simply create a new pool and start scheduling jobs to it.
```go
package main

import "github.com/ernilsson/gopewl"

func main() {
    pool, err := gopewl.NewPool(gopewl.PoolOpts{
	    poolSize: 4,
	    queueSize: 0,
    })  
    if err != nil {
        panic(err)
    }
    defer pool.Close()
    pool.Schedule(func (){
        fmt.Println("Job is being processed")
    })
}
```

To notify the caller when a set of tasks have been completed, use the built-in Go `sync.WaitGroup{}` as shown below.
```go
package main

import "github.com/ernilsson/gopewl"

func main() {
    pool, err := gopewl.NewPool(gopewl.PoolOpts{
	    poolSize: 4,
	    queueSize: 0,
    })
    if err != nil {
        panic(err)
    }
    defer pool.Close()
    wg := sync.WaitGroup{}
    for i := 0; i < 10; i++ {
        wg.Add(1)
        pool.Schedule(func (){
            defer wg.Done()
            // Perform async processing here
        })
    }
    wg.Wait()
    // Continue sync processing here
}
```

### Scheduling with delay
To schedule a job that should run after a given duration, use the `ScheduleWithDelay()` receiver function instead. It 
accepts an additional parameter `delay` and spawns an extra goroutine per call that handles the delaying. The following example
schedules to job to run one second in the future.

```go
package main

import (
	"github.com/ernilsson/gopewl"
	"time"
)

func main() {
	pool, err := gopewl.NewPool(gopewl.PoolOpts{
		poolSize:  4,
		queueSize: 0,
	})
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	for i := 0; i < 10; i++ {
		pool.ScheduleWithDelay(func() {
			// Perform async processing here
		}, time.Second)
	}
	// Continue sync processing here
}
```

## API
### Pool 
The pool type has two receiver methods, both of which are explained below.
#### Schedule
Schedule takes a `func()` and sends it over a channel to any listening worker routine.
#### Close
Close shuts down the job channel and thereby also kills the worker routines that are not currently executing a job. Any
routine that is occupied will exit as soon as the job it is processing has been completed. 
