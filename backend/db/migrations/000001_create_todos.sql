-- +goose Up
-- Crea la tabla de todos.
CREATE TABLE todos (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title      text NOT NULL,
    completed  boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT title_not_blank CHECK (btrim(title) <> '')
);

-- +goose Down
DROP TABLE IF EXISTS todos;
