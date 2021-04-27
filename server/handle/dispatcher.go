package handle

import "strconv"

// Dispatcher keeps the thing together, a controller get job from queue the push to worker
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool		chan chan Job
	MaxWorkers		int
	Workers			[]Worker
	JobQueue		chan Job // A buffered channel that we can send work request on
}

// NewDispatcher create a new instance
func NewDispatcher(maxWorkers, maxQueue int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)
	jobQueue := make(chan Job, maxQueue)
	return &Dispatcher{
		WorkerPool:	workerPool,
		MaxWorkers:	maxWorkers,
		JobQueue: jobQueue,
	}
}


// Start select job from JobQueue push to WorkerPool
func (d *Dispatcher) dispatch(){
	for {
		select {
		case job := <- d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available
				// this will block until a worker is idle
				jobChannel := <- d.WorkerPool

				// dispatch the job to worker job channel that free
				jobChannel <- job
			}(job)
		}
	}
}


// Run start listening and hire the workers
func (d *Dispatcher) Run(){
	// Starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.ID = "worker-" + strconv.Itoa(i + 1)
		d.Workers = append(d.Workers, worker)
		worker.Start()
	}

	// Start select job from JobQueue push to WorkerPool
	go d.dispatch()
}

// DispatchJob send a job to the workers
func (d *Dispatcher) DispatchJob(j *Job){
	go func() {
		// Push the worker onto the queue
		d.JobQueue <- *j
	}()
}


// Stop. Why you ever want to stop?
func (d *Dispatcher) Stop(){
	for _, w := range d.Workers {
		w.Stop()
	}
}