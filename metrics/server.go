package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type BaseStatusResponse struct {
	Obj map[string]int `json:"object"`
}

// Health returns 'READY.'. Used for healthchecking of the service.
func (m *Metrics) Health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "READY.")
}

// Status returns a json dump of BaseStatusResponse
func (m *Metrics) Status(w http.ResponseWriter, r *http.Request) {
	bsr, err := json.Marshal(m.GetAll())
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	w.Write(bsr)
}
