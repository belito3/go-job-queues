package handle

import (
	"encoding/json"
	"io"
	"net/http"
)

type PayloadHandle struct {
}

func NewPayloadHandle() *PayloadHandle {
	return &PayloadHandle{}
}

// Payload the data to be processed, in this example
type Payload struct {
	Magic	string		`json:"magic"`
}

func(h *PayloadHandle) ProcessData(d *Dispatcher, w http.ResponseWriter, r *http.Request) {
	var content Payload
	defer r.Body.Close()

	err := json.NewDecoder(io.LimitReader(r.Body, 2048)).Decode(&content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d.DispatchJob(&Job{content})
	w.WriteHeader(http.StatusOK)
}