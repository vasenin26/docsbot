package metrics

import "testing"

func TestMetricsRegistered(t *testing.T) {
	if RequestsTotal == nil {
		t.Fatal("RequestsTotal is nil")
	}
}
