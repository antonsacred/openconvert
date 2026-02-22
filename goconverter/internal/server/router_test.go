package server

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthEndpoint(t *testing.T) {
	router := newTestRouter()

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
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/v1/conversions", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var body struct {
		Formats map[string][]string `json:"formats"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}

	if len(body.Formats) == 0 {
		t.Fatalf("expected at least one source format")
	}

	if len(body.Formats) != 3 {
		t.Fatalf("expected 3 source formats, got %d (%v)", len(body.Formats), body.Formats)
	}

	targets, ok := body.Formats["png"]
	if !ok {
		t.Fatalf("expected key \"png\" in response, got %v", body.Formats)
	}

	if len(targets) != 2 {
		t.Fatalf("expected exactly two targets for png, got %d (%v)", len(targets), targets)
	}

	if !slices.Contains(targets, "jpg") || !slices.Contains(targets, "webp") {
		t.Fatalf("expected png targets to include [jpg webp], got %v", targets)
	}
}

func TestConvertEndpoint(t *testing.T) {
	router := newTestRouter()

	payload := map[string]string{
		"from":          "png",
		"to":            "jpg",
		"fileName":      "input.png",
		"contentBase64": base64.StdEncoding.EncodeToString(mustEncodePNG(t)),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/convert", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	var response struct {
		From          string `json:"from"`
		To            string `json:"to"`
		FileName      string `json:"fileName"`
		MimeType      string `json:"mimeType"`
		ContentBase64 string `json:"contentBase64"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if response.From != "png" || response.To != "jpg" {
		t.Fatalf("expected from/to png->jpg, got %s->%s", response.From, response.To)
	}
	if response.FileName != "input.jpg" {
		t.Fatalf("expected output fileName input.jpg, got %q", response.FileName)
	}
	if response.MimeType != "image/jpeg" {
		t.Fatalf("expected mime image/jpeg, got %q", response.MimeType)
	}
	if response.ContentBase64 == "" {
		t.Fatalf("expected non-empty contentBase64")
	}

	decoded, err := base64.StdEncoding.DecodeString(response.ContentBase64)
	if err != nil {
		t.Fatalf("failed to decode output base64: %v", err)
	}
	if len(decoded) == 0 {
		t.Fatalf("expected non-empty converted bytes")
	}

	if w.Header().Get("X-Request-Id") == "" {
		t.Fatalf("expected X-Request-Id response header")
	}
}

func TestConvertEndpointRejectsInvalidBase64(t *testing.T) {
	router := newTestRouter()

	payload := `{"from":"png","to":"jpg","fileName":"input.png","contentBase64":"%%%"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/convert", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response struct {
		Error struct {
			Code      string `json:"code"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if response.Error.Code != "invalid_base64" {
		t.Fatalf("expected error code invalid_base64, got %q", response.Error.Code)
	}
}

func TestConvertEndpointRejectsMissingFields(t *testing.T) {
	router := newTestRouter()

	payload := `{"from":"png","to":"jpg","contentBase64":"abcd"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/convert", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response struct {
		Error struct {
			Code      string `json:"code"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if response.Error.Code != "invalid_request" {
		t.Fatalf("expected error code invalid_request, got %q", response.Error.Code)
	}
}

func TestConvertEndpointRejectsUnsupportedPair(t *testing.T) {
	router := newTestRouter()

	payload := map[string]string{
		"from":          "png",
		"to":            "pdf",
		"fileName":      "input.png",
		"contentBase64": base64.StdEncoding.EncodeToString(mustEncodePNG(t)),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/convert", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("expected status %d, got %d", http.StatusUnsupportedMediaType, w.Code)
	}

	var response struct {
		Error struct {
			Code      string `json:"code"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if response.Error.Code != "unsupported_conversion_pair" {
		t.Fatalf("expected error code unsupported_conversion_pair, got %q", response.Error.Code)
	}
}

func TestConvertEndpointRejectsTooLargePayload(t *testing.T) {
	router := newTestRouter()

	oldMax := maxDecodedFileSizeBytes
	maxDecodedFileSizeBytes = 4
	t.Cleanup(func() {
		maxDecodedFileSizeBytes = oldMax
	})

	payload := map[string]string{
		"from":          "png",
		"to":            "jpg",
		"fileName":      "input.png",
		"contentBase64": base64.StdEncoding.EncodeToString([]byte{1, 2, 3, 4, 5}),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/convert", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, w.Code)
	}

	var response struct {
		Error struct {
			Code      string `json:"code"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if response.Error.Code != "payload_too_large" {
		t.Fatalf("expected error code payload_too_large, got %q", response.Error.Code)
	}

	if response.Error.RequestID == "" {
		t.Fatalf("expected requestId in error response")
	}
}

func TestConvertEndpointRejectsTooLargeRequestBody(t *testing.T) {
	router := newTestRouter()

	oldMaxRequestBodyBytes := maxRequestBodyBytes
	maxRequestBodyBytes = 128
	t.Cleanup(func() {
		maxRequestBodyBytes = oldMaxRequestBodyBytes
	})

	payload := map[string]string{
		"from":          "png",
		"to":            "jpg",
		"fileName":      "input.png",
		"contentBase64": strings.Repeat("a", 512),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/convert", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d, got %d", http.StatusRequestEntityTooLarge, w.Code)
	}

	var response struct {
		Error struct {
			Code      string `json:"code"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if response.Error.Code != "payload_too_large" {
		t.Fatalf("expected error code payload_too_large, got %q", response.Error.Code)
	}
}

func TestConvertEndpointReturnsBusyWhenAllSlotsAreUsed(t *testing.T) {
	router := newTestRouter()

	oldMaxConcurrentConversions := maxConcurrentConversions
	maxConcurrentConversions = 0
	t.Cleanup(func() {
		maxConcurrentConversions = oldMaxConcurrentConversions
	})

	conversionConcurrencyMu.Lock()
	currentConcurrentConversions = 0
	conversionConcurrencyMu.Unlock()

	payload := map[string]string{
		"from":          "png",
		"to":            "jpg",
		"fileName":      "input.png",
		"contentBase64": base64.StdEncoding.EncodeToString(mustEncodePNG(t)),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/convert", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
	}

	var response struct {
		Error struct {
			Code      string `json:"code"`
			RequestID string `json:"requestId"`
		} `json:"error"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if response.Error.Code != "converter_busy" {
		t.Fatalf("expected error code converter_busy, got %q", response.Error.Code)
	}
	if response.Error.RequestID == "" {
		t.Fatalf("expected requestId in error response")
	}
}

func TestRequestIDHeaderIsPropagatedToResponse(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("X-Request-Id", "req-test-001")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if got := w.Header().Get("X-Request-Id"); got != "req-test-001" {
		t.Fatalf("expected X-Request-Id to be propagated, got %q", got)
	}
}

func TestOpenAPISpecEndpoint(t *testing.T) {
	router := newTestRouter()

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

	body := w.Body.String()
	if !strings.Contains(body, "/v1/conversions") {
		t.Fatalf("expected OpenAPI body to describe /v1/conversions endpoint")
	}
	if !strings.Contains(body, "/v1/convert") {
		t.Fatalf("expected OpenAPI body to describe /v1/convert endpoint")
	}
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter()
}

func mustEncodePNG(t *testing.T) []byte {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	img.Set(0, 1, color.RGBA{B: 255, A: 255})
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode png fixture: %v", err)
	}

	return buf.Bytes()
}
