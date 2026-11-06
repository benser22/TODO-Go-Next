package handler

import (
	"context"
	"log/slog"
	"net/http"
)

// Pinger is what the readiness probe needs from a dependency. With a
// consumer-side interface, the router only depends on "something it can
// ping", and tests can pass a fake.
type Pinger interface {
	Ping(ctx context.Context) error
}

// handleLiveness answers whether the process is up. It never checks
// dependencies: an outage must not make the orchestrator restart replicas.
func handleLiveness(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "alive"})
}

// healthHandler answers readiness: should the server receive traffic?
type healthHandler struct {
	db Pinger
}

func newHealthHandler(db Pinger) *healthHandler {
	return &healthHandler{db: db}
}

func (h *healthHandler) handleReadiness(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(r.Context()); err != nil {
		slog.WarnContext(r.Context(), "readiness check failed", slog.Any("error", err))
		respondError(w, http.StatusServiceUnavailable, "UNAVAILABLE", "database unavailable")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
