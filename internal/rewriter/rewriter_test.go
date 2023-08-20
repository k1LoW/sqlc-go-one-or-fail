package rewriter

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tenntenn/golden"
)

func TestRewrite(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"default", "testdata/default.query.sql.go"},
		{"emit_methods_with_db_argument: true", "testdata/emit_methods_with_db_argument.query.sql.go"},
		{"emit_result_struct_pointers: true", "testdata/emit_result_struct_pointers.query.sql.go"},
	}
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(bytes.Buffer)
			if err := Rewrite(ctx, tt.in, got); err != nil {
				t.Fatal(err)
			}

			f := filepath.Base(tt.in)
			if os.Getenv("UPDATE_GOLDEN") != "" {
				golden.Update(t, "testdata", f, got)
				return
			}

			if diff := golden.Diff(t, "testdata", f, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
