package server

import "encoding/base64"

type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}

type ConversionsResponse struct {
	Formats map[string][]string `json:"formats"`
}

type ConvertRequest struct {
	From          string `json:"from" example:"png"`
	To            string `json:"to" example:"jpg"`
	FileName      string `json:"fileName" example:"input.png"`
	ContentBase64 string `json:"contentBase64"`
}

type ConvertResponse struct {
	From          string `json:"from" example:"png"`
	To            string `json:"to" example:"jpg"`
	FileName      string `json:"fileName" example:"input.jpg"`
	MimeType      string `json:"mimeType" example:"image/jpeg"`
	ContentBase64 string `json:"contentBase64"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string `json:"code" example:"unsupported_conversion_pair"`
	Message   string `json:"message" example:"conversion from png to pdf is not supported"`
	RequestID string `json:"requestId,omitempty" example:"req-abc123"`
}

var maxDecodedFileSizeBytes = 50 * 1024 * 1024
var maxRequestBodyBytes = int64(base64.StdEncoding.EncodedLen(maxDecodedFileSizeBytes) + (2 * 1024 * 1024))
var maxConcurrentConversions = 4
