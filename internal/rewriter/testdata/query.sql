-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ?;

-- name: GetAuthorName :one
SELECT name FROM authors
WHERE id = ?;

-- name: GetAuthorNameAndBio :one
SELECT name, bio FROM authors
WHERE id = ?;

-- name: GetAuthorBio :one
SELECT bio FROM authors
WHERE id = ?;

-- name: GetAuthorByName :one
SELECT * FROM authors
WHERE name = ?;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
