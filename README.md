# sqlc-go-one-or-fail [![CI](https://github.com/k1LoW/sqlc-go-one-or-fail/actions/workflows/ci.yml/badge.svg)](https://github.com/k1LoW/sqlc-go-one-or-fail/actions/workflows/ci.yml) ![Coverage](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/sqlc-go-one-or-fail/coverage.svg) ![Code to Test Ratio](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/sqlc-go-one-or-fail/ratio.svg) ![Test Execution Time](https://raw.githubusercontent.com/k1LoW/octocovs/main/badges/k1LoW/sqlc-go-one-or-fail/time.svg)

`sqlc-go-one-or-fail` modifies the Go code generated by [sqlc](https://sqlc.dev/) to fail if more than one record is retrieved in the [`:one`](https://docs.sqlc.dev/en/latest/reference/query-annotations.html#one) command.

## Usage

```console
$ sqlc generate
$ sqlc-go-one-or-fail path/to/generated_by_sqlc/*.go
```

### Before

```go
const getAuthorByName = `-- name: GetAuthorByName :one
SELECT id, name, bio FROM authors
WHERE name = ?
`

func (q *Queries) GetAuthorByName(ctx context.Context, name string) (Author, error) {
	row := q.db.QueryRowContext(ctx, getAuthorByName, name)
	var i Author
	err := row.Scan(&i.ID, &i.Name, &i.Bio)
	return i, err
}
```

### After

```go
const getAuthorByName = `-- name: GetAuthorByName :one
SELECT id, name, bio FROM authors
WHERE name = ?
`

func (q *Queries) GetAuthorByName(ctx context.Context, name string) (Author, error) {
	rows, err := q.db.QueryContext(ctx, getAuthorByName, name)
	if err != nil {
		return Author{}, err
	}
	defer rows.Close()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return Author{}, err
		}
		return Author{}, sql.ErrNoRows
	}
	var i Author
	if err := rows.Scan(&i.ID, &i.Name, &i.Bio); err != nil {
		return Author{}, err
	}
	if rows.Next() {
		return Author{}, fmt.Errorf("multiple records were retrieved when the following query was executed: %q", getAuthorByName)
	}
	if err := rows.Close(); err != nil {
		return Author{}, err
	}
	if err := rows.Err(); err != nil {
		return Author{}, err
	}
	return i, err

}
```

## Install

**deb:**

``` console
$ export SQLC_GO_ONE_OR_FAIL_VERSION=X.X.X
$ curl -o sqlc-go-one-or-fail.deb -L https://github.com/k1LoW/sqlc-go-one-or-fail/releases/download/v$SQLC_GO_ONE_OR_FAIL_VERSION/sqlc-go-one-or-fail_$SQLC_GO_ONE_OR_FAIL_VERSION-1_amd64.deb
$ dpkg -i sqlc-go-one-or-fail.deb
```

**RPM:**

``` console
$ export SQLC_GO_ONE_OR_FAIL_VERSION=X.X.X
$ yum install https://github.com/k1LoW/sqlc-go-one-or-fail/releases/download/v$SQLC_GO_ONE_OR_FAIL_VERSION/sqlc-go-one-or-fail_$SQLC_GO_ONE_OR_FAIL_VERSION-1_amd64.rpm
```

**apk:**

``` console
$ export SQLC_GO_ONE_OR_FAIL_VERSION=X.X.X
$ curl -o sqlc-go-one-or-fail.apk -L https://github.com/k1LoW/sqlc-go-one-or-fail/releases/download/v$SQLC_GO_ONE_OR_FAIL_VERSION/runn_$SQLC_GO_ONE_OR_FAIL_VERSION-1_amd64.apk
$ apk add sqlc-go-one-or-fail.apk
```

**homebrew tap:**

```console
$ brew install k1LoW/tap/sqlc-go-one-or-fail
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/sqlc-go-one-or-fail/releases)

**go install:**

```console
$ go install github.com/k1LoW/sqlc-go-one-or-fail@latest
```
