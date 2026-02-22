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
		PossibleConvertationFormats map[string][]string `json:"possibleConvertationFormats"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if len(body.PossibleConvertationFormats) == 0 {
		t.Fatalf("expected at least one source format")
	}

	targets, ok := body.PossibleConvertationFormats["png"]
	if !ok {
		t.Fatalf("expected key \"png\" in response, got %v", body.PossibleConvertationFormats)
	}

	if len(targets) != 1 {
		t.Fatalf("expected exactly one target for png, got %d (%v)", len(targets), targets)
	}

	if targets[0] != "jpg" {
		t.Fatalf("expected png target to be jpg, got %v", targets)
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
