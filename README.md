# prioritized
A simple work priority implementation in go.

## Why
I was stress testing a ws server, it turned out can't properly serve clients connected before clients connected later.

This simple implementation allows users to queue jobs into high/low priority channel. Jobs in high priority channels will always be dequeued before low priority jobs.

This kind of manual prioritizing is needed because goroutines are automatically scheduled by go and we have no way to control which goroutine get awoken when the total goroutines exceeds GOMAXPROCS.

## Usage
The implementation is quite simple. You just call `prioritized.Init(jobQueueCapacity int, workerConcurrency int, closeWorkerSignal chan struct{})` to get a `WorkerGroup`, then call define job types, then stuff in your job to either high or low priority queue.

#### Init
```go
prioritized.Init(jobQueueCapacity int, workerConcurrency int, closeWorkerSignal chan struct{})
```
- **jobQueueCapacity**: maximum amount of jobs in the both the high and low priority queue, the channel will be blocked while it's full.
- **jobQueueConcurrency**: How many workers will be created, that is, how many workers will be concurrently working on queued jobs.
- **closeWorkerSignal**: close the channel to close all the workers.

#### Defining Jobs
```go
workerGroup.DefineJobType(typeSymbol interface{}, jobHandler func(jobData interface{}))
```
- **typeSymbol**: could be anything, used as key to call corresponding job handlers.
- **jobHandler**: the func to handle the job.

#### Pending Jobs
##### Blocking
```go
workerGroup.PendHighPriorityJob(typeSymbol interface{}, jobData interface{})
workerGroup.PendLowPriorityJob(typeSymbol interface{}, jobData interface{})
workerGroup.PendAutoFallback(typeSymbol interface{}, jobData interface{}) // Try HP queue first, if it's full, go to LP
```
- **typeSymbol**: the key corresponding the desired job handler.
- **jobData**: the data to be handled.

##### Non-Blocking
These methods does not block if the queues are full, instead it tells you if it's successful.
```go
workerGroup.TryPendHighPriorityJo(typeSymbol interface{}, jobData interface{}) bool
workerGroup.TryPendLowPriorityJo(typeSymbol interface{}, jobData interface{}) bool
workerGroup.TryPendAutoFallback(typeSymbol interface{}, jobData interface{}) bool
```
