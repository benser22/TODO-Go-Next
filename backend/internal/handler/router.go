package handler

import (
	"net/http"

	"github.com/benja/todo-api/internal/service"
)

// NewRouter builds the HTTP handler of the API: routes first, then
// middlewares composed outside-in (recoverer must be the outermost one).
func NewRouter(allowedOrigin string, todoSvc *service.TodoService, pinger Pinger) http.Handler {
	h := NewTodoHandler(todoSvc)

	mux := http.NewServeMux()

	// Health probes. Liveness stays dependency-free by design.
	mux.HandleFunc("GET /healthz", handleLiveness)
	mux.HandleFunc("GET /readyz", newHealthHandler(pinger).handleReadiness)

	// Todos.
	mux.HandleFunc("GET /api/v1/todos", h.list)
	mux.HandleFunc("POST /api/v1/todos", h.create)
	mux.HandleFunc("GET /api/v1/todos/{id}", h.get)
	mux.HandleFunc("PATCH /api/v1/todos/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/todos/{id}", h.delete)

	return recoverer(cors(allowedOrigin, requestID(requestLogger(mux))))
}
