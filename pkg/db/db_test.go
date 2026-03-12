package db

import (
	"context"
	"testing"

	_ "modernc.org/sqlite"
)

func TestNewInMemory(t *testing.T) {
	d, err := New("sqlite", "file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	ctx := context.Background()
	_, err = d.Exec(ctx, "CREATE TABLE t (id INTEGER PRIMARY KEY)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = d.Exec(ctx, "INSERT INTO t (id) VALUES (1)")
	if err != nil {
		t.Fatal(err)
	}
	var n int
	row := d.QueryRow(ctx, "SELECT id FROM t WHERE id = 1")
	if err := row.Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("want 1 got %d", n)
	}
}
