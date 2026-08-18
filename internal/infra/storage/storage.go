// Package storage provides a generic object-storage boundary. Real providers
// (S3, GCS, etc.) should implement the Store interface. This template ships
// with a local filesystem implementation.
package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Store is the generic storage interface modules can depend on.
type Store interface {
	// Put writes src to the object at key.
	Put(ctx context.Context, key string, src io.Reader) error
	// Get opens the object at key for reading.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete removes the object at key.
	Delete(ctx context.Context, key string) error
}

// LocalStore stores objects on the local filesystem.
type LocalStore struct {
	baseDir string
}

// NewLocalStore creates a LocalStore rooted at baseDir.
func NewLocalStore(baseDir string) (*LocalStore, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("storage: create base dir: %w", err)
	}
	return &LocalStore{baseDir: baseDir}, nil
}

// Put writes src to baseDir/key.
func (s *LocalStore) Put(ctx context.Context, key string, src io.Reader) error {
	path := filepath.Join(s.baseDir, key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("storage: mkdir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("storage: create: %w", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, src); err != nil {
		return fmt.Errorf("storage: write: %w", err)
	}
	return nil
}

// Get opens baseDir/key.
func (s *LocalStore) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(s.baseDir, key))
	if err != nil {
		return nil, fmt.Errorf("storage: open: %w", err)
	}
	return f, nil
}

// Delete removes baseDir/key.
func (s *LocalStore) Delete(ctx context.Context, key string) error {
	if err := os.Remove(filepath.Join(s.baseDir, key)); err != nil {
		return fmt.Errorf("storage: delete: %w", err)
	}
	return nil
}
