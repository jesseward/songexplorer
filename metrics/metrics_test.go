package metrics

import (
	"testing"
)

// TestMetricsIncAndGet ensures increment and fetch functionality works as expected.
func TestMetricsIncAndGet(t *testing.T) {

	tt := []struct {
		value int64
	}{
		{1},
		{2},
		{3},
	}

	m := New()

	for _, tc := range tt {
		m.Inc("test-counter")
		c := m.Get("test-counter")
		if c != tc.value {
			t.Errorf("expected %d , received %d", tc.value, c)
		}
	}
}

// TestMetricsGetAll verifies results of all metrics in counter.
func TestMetricsGetAll(t *testing.T) {
	m := New()
	m.Inc("test-counter-a")
	m.Inc("test-counter-b")
	m.Inc("test-counter-a")
	v := m.GetAll()

	if v["test-counter-a"] != 2 {
		t.Errorf("exptected 2, received %d", v["test-counter-a"])
	}

	if v["test-counter-b"] != 1 {
		t.Errorf("exptected 1, received %d", v["test-counter-a"])
	}
}
