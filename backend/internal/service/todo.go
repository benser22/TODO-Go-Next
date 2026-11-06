// Package service contains the use cases of the application. It orchestrates
// the domain and defines the interfaces it needs from the persistence layer,
// so any store implementation satisfies them implicitly.
package service

import (
	"context"
	"fmt"

	"github.com/benja/todo-api/internal/domain"
)

// TodoStore is what TodoService needs from a persistence layer.
// The interface is defined at the consumer, not at the implementation.
type TodoStore interface {
	List(ctx context.Context) ([]domain.Todo, error)
	Get(ctx context.Context, id string) (domain.Todo, error)
	Create(ctx context.Context, title string) (domain.Todo, error)
	Update(ctx context.Context, t domain.Todo) (domain.Todo, error)
	Delete(ctx context.Context, id string) error
}

// TodoService implements the todo use cases on top of a TodoStore.
type TodoService struct {
	store TodoStore
}

// NewTodoService returns a ready-to-use service.
func NewTodoService(store TodoStore) *TodoService {
	return &TodoService{store: store}
}

// List returns every todo, newest first.
func (s *TodoService) List(ctx context.Context) ([]domain.Todo, error) {
	todos, err := s.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	return todos, nil
}

// Get returns a single todo by ID.
func (s *TodoService) Get(ctx context.Context, id string) (domain.Todo, error) {
	todo, err := s.store.Get(ctx, id)
	if err != nil {
		return domain.Todo{}, fmt.Errorf("get todo: %w", err)
	}
	return todo, nil
}

// Create validates the title and persists a new todo.
func (s *TodoService) Create(ctx context.Context, title string) (domain.Todo, error) {
	title, err := domain.ValidateTitle(title)
	if err != nil {
		return domain.Todo{}, err
	}
	todo, err := s.store.Create(ctx, title)
	if err != nil {
		return domain.Todo{}, fmt.Errorf("create todo: %w", err)
	}
	return todo, nil
}

// Update replaces the mutable fields of a todo after validating its title.
func (s *TodoService) Update(ctx context.Context, t domain.Todo) (domain.Todo, error) {
	title, err := domain.ValidateTitle(t.Title)
	if err != nil {
		return domain.Todo{}, err
	}
	t.Title = title
	todo, err := s.store.Update(ctx, t)
	if err != nil {
		return domain.Todo{}, fmt.Errorf("update todo: %w", err)
	}
	return todo, nil
}

// Delete removes a todo by ID.
func (s *TodoService) Delete(ctx context.Context, id string) error {
	if err := s.store.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	return nil
}
