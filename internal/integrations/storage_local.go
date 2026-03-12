package integrations

import (
	"fmt"
	"os"
	"path/filepath"
)

// LocalStorage stores files on disk.
type LocalStorage struct {
	root string
}

// NewLocalStorage creates a local storage at root path.
func NewLocalStorage(root string) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{root: root}, nil
}

// Presign returns a local file path (no real presign for local).
func (l *LocalStorage) Presign(bucket, key string, _ int) (string, error) {
	return filepath.Join(l.root, bucket, key), nil
}

// Put writes data to root/bucket/key.
func (l *LocalStorage) Put(bucket, key string, data []byte) error {
	p := filepath.Join(l.root, bucket, key)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

// Get reads file at root/bucket/key.
func (l *LocalStorage) Get(bucket, key string) ([]byte, error) {
	p := filepath.Join(l.root, bucket, key)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", p)
	}
	return os.ReadFile(p)
}
