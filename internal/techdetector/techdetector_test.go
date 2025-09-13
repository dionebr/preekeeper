package techdetector

import (
	"testing"
)

func TestEngine_Fingerprint_Smoke(t *testing.T) {
	eng, err := New()
	if err != nil {
		t.Skipf("engine init failed, skipping integration test: %v", err)
	}

	headers := map[string][]string{"Server": {"nginx"}}
	body := []byte("<html><head><meta name=\"generator\" content=\"Test\"></head><body></body></html>")
	res := eng.Fingerprint(headers, body)
	if res == nil {
		t.Fatalf("expected non-nil result from Fingerprint")
	}
}
