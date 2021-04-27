package handle

import (
	"github.com/gorilla/websocket"
	"job_queue/pkg/logger"
	"net/http"
)

// Use a mutex for production and GorillaWebsocket
var wsList [] *websocket.Conn

// WorkerUpdate worker status update  msg to clients
type WorkerUpdate struct {
	WorkerID	string 	`json:"WorkerID"`
	Status		string	`json:"Status"`
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	wsList = append(wsList, ws)
}

// clientStream send status update for all workers to all clients
func clientStream(w *Worker, status string) {
	m := WorkerUpdate{w.ID, status}

	for _, ws := range wsList {
		err :=  ws.WriteJSON(m)

		if err != nil {
			logger.Errorf(nil, err.Error())
		}
	}
}