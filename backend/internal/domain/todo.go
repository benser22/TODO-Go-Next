// Package domain holds the core Todo entity and its business rules.
// It must not import infrastructure packages (HTTP, database, handlers).
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// MaxTitleLength is the upper bound for a todo title.
const MaxTitleLength = 200

// Todo is the core entity of the application.
type Todo struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Errors returned by domain rules. Callers register them with errors.Is.
var (
	ErrNotFound     = errors.New("todo: not found")
	ErrInvalidID    = errors.New("todo: invalid id")
	ErrInvalidTitle = errors.New("todo: title must not be empty")
)

// ValidateTitle normalizes a title (trims surrounding spaces) and enforces
// its constraints.
func ValidateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", ErrInvalidTitle
	}
	if len(title) > MaxTitleLength {
		return "", fmt.Errorf("%w: exceeds %d characters", ErrInvalidTitle, MaxTitleLength)
	}
	return title, nil
}
