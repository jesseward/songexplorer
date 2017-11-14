package metrics

import (
	"net/http"
	"sync"
)

type Metrics struct {
	mu      sync.RWMutex
	counter map[string]int64
}

func (m *Metrics) GetAll() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}

func (m *Metrics) Inc(k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[k]++
}

func New() *Metrics {
	return &Metrics{counter: make(map[string]int64)}
}

func NewAdminServletMetrics() (*Metrics, *http.ServeMux) {
	m := New()
	mux := http.NewServeMux()
	mux.HandleFunc("/health/", m.Health)
	mux.HandleFunc("/status/", m.Status)
	return m, mux
}
