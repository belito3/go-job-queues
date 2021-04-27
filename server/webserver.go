package server

import (
	"github.com/gin-gonic/gin"
	"job_queue/server/config"
	"job_queue/server/handle"
	"job_queue/server/route"
)

// StartServer start an web server, which can be used to control the jobs
func StartServer(d *handle.Dispatcher) {
	port := config.AppConf.ServicePort
	app := gin.Default()

	route.ApplyRoutes(app, d)

	_ = app.Run(port)
}
