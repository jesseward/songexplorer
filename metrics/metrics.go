package metrics

import (
	"net/http"
	"sync"
)

// Metrics container for stats/counter
type Metrics struct {
	mu      sync.RWMutex
	counter map[string]int64
}

// GetAll retrievs all metrics and their state.
func (m *Metrics) GetAll() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.counter
}

// Get retrieves the value of a single key
func (m *Metrics) Get(k string) int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if v, ok := m.counter[k]; ok {
		return v
	}
	return 0
}

// Inc +1's to a key.
func (m *Metrics) Inc(k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counter[k]++
}

// New creates a new metrics object.
func New() *Metrics {
	return &Metrics{counter: make(map[string]int64)}
}

// NewAdminServletMetrics is the default servlet for metrics daemon.
func NewAdminServletMetrics() (*Metrics, *http.ServeMux) {
	m := New()
	mux := http.NewServeMux()
	mux.HandleFunc("/health/", m.Health)
	mux.HandleFunc("/status/", m.Status)
	return m, mux
}
