package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

// ctxKeyType is the unexported context key type; string keys can be
// overwritten by any other package in the process.
type ctxKeyType struct{}

// requestIDFrom reads the request identifier the requestID middleware
// attached to the request context.
func requestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyType{}).(string)
	return id
}

// requestID assigns an identifier per request: the incoming X-Request-ID
// when present, a random one otherwise. It is echoed back to the client.
func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			var buf [8]byte
			_, _ = rand.Read(buf[:])
			id = hex.EncodeToString(buf[:])
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyType{}, id)))
	})
}

// recoverer must wrap everything else: a panic in a handler is converted
// into a 500 response instead of a dead process.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("panic recovered",
					slog.Any("panic", rec),
					slog.String("stack", string(debug.Stack())),
					slog.String("request_id", requestIDFrom(r.Context())),
				)
				respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// statusRecorder captures the response status line for request logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// requestLogger writes one log line per request with timing and status.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Duration("duration", time.Since(start).Round(time.Millisecond)),
			slog.String("request_id", requestIDFrom(r.Context())),
		}
		switch {
		case rec.status >= http.StatusInternalServerError:
			slog.LogAttrs(r.Context(), slog.LevelError, "request", attrs...)
		case rec.status >= http.StatusBadRequest:
			slog.LogAttrs(r.Context(), slog.LevelWarn, "request", attrs...)
		default:
			slog.LogAttrs(r.Context(), slog.LevelInfo, "request", attrs...)
		}
	})
}

// cors allows the configured web origin to call the API and answers
// preflight requests directly.
func cors(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Add("Vary", "Origin")
		}
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "300")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
