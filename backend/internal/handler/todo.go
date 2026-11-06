package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/benja/todo-api/internal/domain"
	"github.com/benja/todo-api/internal/service"
)

// TodoHandler serves the /api/v1/todos endpoints. It maps domain errors to
// status codes at this boundary and never exposes internal error text.
type TodoHandler struct {
	svc *service.TodoService
}

// NewTodoHandler returns a handler that depends only on its service.
func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

// todoDTO is the wire format of a todo: camelCase keys and RFC 3339
// timestamps.
type todoDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toDTO(t domain.Todo) todoDTO {
	return todoDTO{
		ID:        t.ID,
		Title:     t.Title,
		Completed: t.Completed,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

// todoListDTO wraps the list so the wire shape stays extensible.
type todoListDTO struct {
	Todos []todoDTO `json:"todos"`
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type updateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}

// list handles GET /api/v1/todos.
func (h *TodoHandler) list(w http.ResponseWriter, r *http.Request) {
	todos, err := h.svc.List(r.Context())
	if err != nil {
		writeServiceError(w, r, err)
		return
	}

	dt := make([]todoDTO, 0, len(todos))
	for _, t := range todos {
		dt = append(dt, toDTO(t))
	}
	respondJSON(w, http.StatusOK, todoListDTO{Todos: dt})
}

// get handles GET /api/v1/todos/{id}.
func (h *TodoHandler) get(w http.ResponseWriter, r *http.Request) {
	todo, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if writeServiceError(w, r, err) {
		return
	}
	respondJSON(w, http.StatusOK, toDTO(todo))
}

// create handles POST /api/v1/todos.
func (h *TodoHandler) create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if !decode(w, r, &req) {
		return
	}

	todo, err := h.svc.Create(r.Context(), req.Title)
	if writeServiceError(w, r, err) {
		return
	}
	respondJSON(w, http.StatusCreated, toDTO(todo))
}

// update handles PATCH /api/v1/todos/{id} with a partial payload: omitted
// fields keep their current value.
func (h *TodoHandler) update(w http.ResponseWriter, r *http.Request) {
	var req updateTodoRequest
	if !decode(w, r, &req) {
		return
	}

	current, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if writeServiceError(w, r, err) {
		return
	}
	if req.Title != nil {
		current.Title = *req.Title
	}
	if req.Completed != nil {
		current.Completed = *req.Completed
	}

	updated, err := h.svc.Update(r.Context(), current)
	if writeServiceError(w, r, err) {
		return
	}
	respondJSON(w, http.StatusOK, toDTO(updated))
}

// delete handles DELETE /api/v1/todos/{id}.
func (h *TodoHandler) delete(w http.ResponseWriter, r *http.Request) {
	err := h.svc.Delete(r.Context(), r.PathValue("id"))
	if writeServiceError(w, r, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// writeServiceError maps a service/domain error to its HTTP response and
// reports whether the request was answered: false means err is nil and the
// caller continues on the happy path.
func writeServiceError(w http.ResponseWriter, r *http.Request, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, domain.ErrNotFound):
		respondError(w, http.StatusNotFound, "NOT_FOUND", "todo not found")
	case errors.Is(err, domain.ErrInvalidID):
		respondError(w, http.StatusBadRequest, "INVALID_ID", "id must be a valid uuid")
	case errors.Is(err, domain.ErrInvalidTitle):
		respondError(w, http.StatusBadRequest, "INVALID_TITLE",
			"title must not be empty and accepts up to "+strconv.Itoa(domain.MaxTitleLength)+" characters")
	default:
		slog.LogAttrs(r.Context(), slog.LevelError, "unhandled service error",
			slog.Any("error", err),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("request_id", requestIDFrom(r.Context())),
		)
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected server error")
	}
	return true
}
