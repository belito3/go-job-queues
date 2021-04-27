package main

import (
	"flag"
	"job_queue/pkg/logger"
	"job_queue/server"
	"job_queue/server/handle"
)

func main() {
	MaxWorker := flag.Uint("MAX_WORKERS", 5, "max nr of workers")
	MaxQueue := flag.Uint("MAX_QUEUE", 10, "max nr of jobs in queue")
	flag.Parse()

	logger.Infof(nil,"open http://localhost:8080/index in your browser & keep this process open.")
	dispatcher := handle.NewDispatcher(int(*MaxWorker), int(*MaxQueue))
	// Run start listening and hire the workers
	dispatcher.Run()
	server.StartServer(dispatcher)
}
