# gopewl 🌊
A simple package allowing quick capping of goroutines for rudimentary asynchronous tasks.

## Example
The API is trivial to use. Simply create a new pool and start scheduling jobs to it.
```go
const numOfWorkers = 4

func main() {
    pool, err := gopewl.NewPool(numOfWorkers, 0)
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
const numOfWorkers = 4

func main() {
    pool, err := gopewl.NewPool(numOfWorkers, 0)
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

## API
### Pool 
The pool type has two receiver methods, both of which are explained below.
#### Schedule
Schedule takes a `func()` and sends it over a channel to any listening worker routine.
#### Close
Close shuts down the job channel and thereby also kills the worker routines that are not currently executing a job. Any
routine that is occupied will exit as soon as the job it is processing has been completed. 
