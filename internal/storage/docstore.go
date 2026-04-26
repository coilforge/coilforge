package storage

import "time"

// DocInfo is one document entry in the single-folder browser.
type DocInfo struct {
	Name       string    // file name shown in the browser list.
	SizeBytes  int64     // file size in bytes.
	ModifiedAt time.Time // file modification time.
}

// DocStore defines simple single-folder document operations.
type DocStore interface {
	ListDocs() ([]DocInfo, error)
	LoadDoc(name string) ([]byte, error)
	SaveDoc(name string, data []byte) error
	DeleteDoc(name string) error
}
