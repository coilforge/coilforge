package storage

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	appDirName  = "coilforge"
	docsDirName = "docs"
)

// LocalFSStore stores project docs in a single local folder.
type LocalFSStore struct {
	rootDir string
}

// NewLocalFSStore creates a LocalFSStore rooted at rootDir.
func NewLocalFSStore(rootDir string) *LocalFSStore {
	return &LocalFSStore{rootDir: rootDir}
}

// NewDefaultLocalFSStore creates a LocalFSStore at the platform config directory.
func NewDefaultLocalFSStore() (*LocalFSStore, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return NewLocalFSStore(filepath.Join(cfgDir, appDirName, docsDirName)), nil
}

func (store *LocalFSStore) ListDocs() ([]DocInfo, error) {
	if err := os.MkdirAll(store.rootDir, 0o755); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(store.rootDir)
	if err != nil {
		return nil, err
	}
	docs := make([]DocInfo, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		docs = append(docs, DocInfo{
			Name:       entry.Name(),
			SizeBytes:  info.Size(),
			ModifiedAt: info.ModTime(),
		})
	}
	sort.Slice(docs, func(i, j int) bool {
		return strings.ToLower(docs[i].Name) < strings.ToLower(docs[j].Name)
	})
	return docs, nil
}

func (store *LocalFSStore) LoadDoc(name string) ([]byte, error) {
	path, err := store.docPath(name)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func (store *LocalFSStore) SaveDoc(name string, data []byte) error {
	path, err := store.docPath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(store.rootDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (store *LocalFSStore) DeleteDoc(name string) error {
	path, err := store.docPath(name)
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func (store *LocalFSStore) docPath(name string) (string, error) {
	clean := strings.TrimSpace(name)
	if clean == "" {
		return "", errors.New("document name is required")
	}
	if strings.Contains(clean, "/") || strings.Contains(clean, "\\") {
		return "", errors.New("document name must not contain path separators")
	}
	return filepath.Join(store.rootDir, clean), nil
}
