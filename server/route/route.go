package route

import (
	"github.com/gin-gonic/gin"
	"job_queue/pkg/logger"
	"job_queue/server/handle"
	"net/http"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.Engine, d *handle.Dispatcher) {
	payloadHandle := handle.NewPayloadHandle()

	// Serving static file
	r.StaticFS("/index", http.Dir("./client"))

	r.GET("/ws", func(c *gin.Context) {
		handle.WebsocketHandler(c.Writer, c.Request)
	})

	// Use gin with standard handlers
	// enclosure so we can send the dispatcher, avoiding global variable
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		logger.Infof(nil,"received /job request %v %v %v")

		payloadHandle.ProcessData(d, w, r)
	}
	r.POST("/job", gin.WrapF(wrapper))
}