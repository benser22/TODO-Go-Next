-- name: ListTodos :many
SELECT id, title, completed, created_at, updated_at
FROM todos
ORDER BY created_at DESC, id DESC;

-- name: GetTodo :one
SELECT id, title, completed, created_at, updated_at
FROM todos
WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todos (title)
VALUES ($1)
RETURNING id, title, completed, created_at, updated_at;

-- name: UpdateTodo :one
UPDATE todos
SET title = $2, completed = $3, updated_at = now()
WHERE id = $1
RETURNING id, title, completed, created_at, updated_at;

-- name: DeleteTodo :execrows
DELETE FROM todos
WHERE id = $1;
