package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"testdata/gen"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	tmpdir, err := os.MkdirTemp(os.TempDir(), "rewriter-test")
	if err != nil {
		return err
	}
	p := filepath.Join(tmpdir, "test.db")
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		return err
	}
	defer db.Close()
	b, err := os.ReadFile("schema.sql")
	if err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, string(b)); err != nil {
		return err
	}

	q := gen.New()
	name := "foo"

	for i := 0; i < 10; i++ {
		if _, err := q.CreateAuthor(ctx, db, gen.CreateAuthorParams{
			Name: name,
		}); err != nil {
			return err
		}
	}
	if _, err := q.GetAuthorByName(ctx, db, name); err != nil {
		return err
	}
	return nil
}
