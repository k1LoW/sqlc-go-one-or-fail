// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package gen

import (
	"context"
	"database/sql"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
)
RETURNING id, name, bio
`

type CreateAuthorParams struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) CreateAuthor(ctx context.Context, arg CreateAuthorParams) (*Author, error) {
	row := q.db.QueryRowContext(ctx, createAuthor, arg.Name, arg.Bio)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return &i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAuthor, id)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = ?
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (*Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return &i, err
}

const getAuthorBio = `-- name: GetAuthorBio :one
SELECT bio FROM authors
WHERE id = ?
`

func (q *Queries) GetAuthorBio(ctx context.Context, id int64) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getAuthorBio, id)
	var bio sql.NullString
	err := row.Scan(&bio)
	return bio, err
}

const getAuthorByName = `-- name: GetAuthorByName :one
SELECT id, name, bio FROM authors
WHERE name = ?
`

func (q *Queries) GetAuthorByName(ctx context.Context, name string) (*Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthorByName, name)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return &i, err
}

const getAuthorName = `-- name: GetAuthorName :one
SELECT name FROM authors
WHERE id = ?
`

func (q *Queries) GetAuthorName(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRowContext(ctx, getAuthorName, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

const getAuthorNameAndBio = `-- name: GetAuthorNameAndBio :one
SELECT name, bio FROM authors
WHERE id = ?
`

type GetAuthorNameAndBioRow struct {
	Name string
	Bio  sql.NullString
}

func (q *Queries) GetAuthorNameAndBio(ctx context.Context, id int64) (*GetAuthorNameAndBioRow, error) {
	row := q.db.QueryRowContext(ctx, getAuthorNameAndBio, id)
	var i GetAuthorNameAndBioRow
	err := row.Scan(&i.Name, &i.Bio)
	return &i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name
`

func (q *Queries) ListAuthors(ctx context.Context) ([]*Author, error) {
	rows, err := q.db.QueryContext(ctx, listAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
