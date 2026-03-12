package integrations

import (
	"path/filepath"
	"testing"
)

func TestLocalStoragePutGet(t *testing.T) {
	dir := t.TempDir()
	s, err := NewLocalStorage(dir)
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("hello")
	if err := s.Put("bucket", "key", data); err != nil {
		t.Fatal(err)
	}
	got, err := s.Get("bucket", "key")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(data) {
		t.Errorf("want %q got %q", data, got)
	}
}

func TestLocalStoragePresign(t *testing.T) {
	dir := t.TempDir()
	s, err := NewLocalStorage(dir)
	if err != nil {
		t.Fatal(err)
	}
	path, err := s.Presign("b", "k", 60)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "b", "k")
	if path != want {
		t.Errorf("Presign want %q got %q", want, path)
	}
}

func TestLocalStorageGetMissing(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewLocalStorage(dir)
	_, err := s.Get("b", "missing")
	if err == nil {
		t.Fatal("Get missing file should return error")
	}
}
