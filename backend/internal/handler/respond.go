package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// apiError is the stable wire format for API errors.
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

// maxBodyBytes caps JSON request bodies (1 MiB) to guard against oversized
// payloads.
const maxBodyBytes = 1 << 20

// respondJSON writes v as JSON. Encoding failures are only logged — the
// status is already committed, so there is nothing to tell the client.
func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", slog.Any("error", err))
	}
}

// respondError writes the JSON error envelope. Internal error text must
// never reach the client; handlers log it separately.
func respondError(w http.ResponseWriter, status int, code, message string) {
	respondJSON(w, status, errorEnvelope{Error: apiError{Code: code, Message: message}})
}

// decode reads a JSON request body into dst with an explicit size cap and
// unknown fields rejected. On failure it answers 400 and returns false.
func decode[T any](w http.ResponseWriter, r *http.Request, dst *T) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_JSON", "request body is not a valid JSON object")
		return false
	}
	return true
}
