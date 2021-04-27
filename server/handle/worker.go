package handle

import (
	"errors"
	"job_queue/pkg/logger"
	"math/rand"
	"time"
)


// Job represents the task/job to run with the payload
type Job struct {
	Payload Payload
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool 			chan chan Job  // contain chan worker free
	JobChannel			chan Job
	quit				chan bool
	dispatcherJobQueue	chan Job
	ID					string // for debugging purposes
}

// NewWorker workers are the foundation of our queue system
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit: 		make(chan bool),
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need stop it
func (w *Worker) Start(){
	go func() {
		for {
			// register the current worker into the worker queue (worker pool)
			// this is the worker's way to say "I'm free! Give me a job!"
			w.WorkerPool <- w.JobChannel

			// clientStream send status update for all workers to all clients
			clientStream(w, "IDLE")

			select {
			case job := <- w.JobChannel:
				// worker have received a job
				clientStream(w, "received: " + job.Payload.Magic)
				time.Sleep(2  * time.Second) // fake time

				clientStream(w, "processing: " + job.Payload.Magic)
				//err := w.process(w, job)
				err := w.processData(job)

				if err == nil {
					clientStream(w, "finished: " + job.Payload.Magic)
				} else {
					clientStream(w, "failed: " + job.Payload.Magic)
				}

				time.Sleep(2 * time.Second)  // fake time
			case <- w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}


// every payload (request to /job) from the client is sent here
func (w *Worker) processData(j Job) error{
	text := j.Payload.Magic
	if text == "error" {
		return errors.New("error")
	}
	logger.Infof(nil,"processing '%v' by %v\n", text, w.ID)
	//simulating a very long time to process
	//so we can understand the process
	time.Sleep(time.Duration(rand.Intn(5)+3) * time.Second)
	logger.Infof(nil,"done processing '%v' by %v\n", text, w.ID)
	return nil
}

// Stop signals the worker to s top listening for worker requests
func (w Worker) Stop(){
	go func() {
		w.quit <- true
	}()
}