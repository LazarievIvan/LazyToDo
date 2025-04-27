-- name: CreateTodo :one
INSERT INTO todos (description, status, created, updated)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTodos :many
SELECT * FROM todos
ORDER BY id;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = $1 LIMIT 1;

-- name: UpdateTodo :one
UPDATE todos
SET description = $2, status = $3, updated = $4
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;