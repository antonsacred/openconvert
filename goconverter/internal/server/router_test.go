package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status field to be 'ok', got %q", body["status"])
	}
}

func TestConversionsEndpoint(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/conversions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body struct {
		Output [][]string `json:"output"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if len(body.Output) == 0 {
		t.Fatalf("expected at least one conversion pair")
	}

	for i, pair := range body.Output {
		if len(pair) != 2 {
			t.Fatalf("expected conversion pair at index %d to have exactly 2 entries, got %d", i, len(pair))
		}
		if pair[0] == "" || pair[1] == "" {
			t.Fatalf("expected non-empty formats in pair at index %d, got %v", i, pair)
		}
	}
}

func TestOpenAPISpecEndpoint(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected Content-Type to include application/json, got %q", contentType)
	}

	if !strings.Contains(w.Body.String(), "/conversions") {
		t.Fatalf("expected OpenAPI body to describe /conversions endpoint")
	}
}
