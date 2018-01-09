package storagedir_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/anothermemory/storage"
	"github.com/anothermemory/storagedir"
	"github.com/anothermemory/storagetests"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDirectoryInMemoryStorage(t *testing.T) {
	storagetests.RunStorageTests(t, createDirectoryInMemoryStorage, nil)
}

func TestDirectoryStorage(t *testing.T) {
	dir, err := ioutil.TempDir("", "storage_directory_root")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	storagetests.RunStorageTests(t, func() storage.Storage {
		return storagedir.NewDirectoryStorage(path.Join(dir, uuid.NewV4().String()))
	}, loadDirectoryStorageFromJSON)
}

func createDirectoryInMemoryStorage() storage.Storage {
	return storagedir.NewDirectoryInMemoryStorage()
}

func loadDirectoryStorageFromJSON(b []byte) (storage.Storage, error) {
	return storagedir.NewDirectoryStorageFromJSONConfig(b)
}
