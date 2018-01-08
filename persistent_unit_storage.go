package storagedir

import (
	"os"

	"github.com/anothermemory/storage"
)

type persistentUnitStorage interface {
	storage.Storage
	mkdirAll(path string, perm os.FileMode) error
	removeDir(name string) error
	writeFile(filename string, data []byte, perm os.FileMode) error
	readFile(filename string) ([]byte, error)
}
