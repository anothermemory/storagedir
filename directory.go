package storagedir

import (
	"encoding/json"
	"os"

	"github.com/anothermemory/storage"
	"github.com/anothermemory/unit"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// TypeDirectory represents storage type which stores all items in directory
const TypeDirectory = "directory"

// TypeDirectoryInMemory represents storage type which stores all items in directory stored in memory
const TypeDirectoryInMemory = "directory_in_memory"

type directoryStorage struct {
	rootDir     string
	fs          afero.Fs
	fsUtil      *afero.Afero
	storageType string
}

// NewDirectoryStorage creates new storage which uses filesystem to store units
func NewDirectoryStorage(rootDir string) storage.Storage {
	fs := afero.NewOsFs()
	return &directoryStorage{rootDir: rootDir, fs: fs, fsUtil: &afero.Afero{Fs: fs}, storageType: TypeDirectory}
}

// NewDirectoryInMemoryStorage creates new storage which uses memory to store units
func NewDirectoryInMemoryStorage() storage.Storage {
	fs := afero.NewMemMapFs()
	return &directoryStorage{rootDir: "/anothermemory", fs: fs, fsUtil: &afero.Afero{Fs: fs}, storageType: TypeDirectoryInMemory}
}

// NewDirectoryStorageFromJSONConfig creates new storage from it's serialized JSON configuration
func NewDirectoryStorageFromJSONConfig(b []byte) (storage.Storage, error) {
	var s directoryStorage
	err := json.Unmarshal(b, &s)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *directoryStorage) mkdirAll(path string, perm os.FileMode) error {
	return s.fs.MkdirAll(path, perm)
}

func (s *directoryStorage) removeDir(name string) error {
	return s.fs.RemoveAll(name)
}

func (s *directoryStorage) writeFile(filename string, data []byte, perm os.FileMode) error {
	return s.fsUtil.WriteFile(filename, data, perm)
}

func (s *directoryStorage) readFile(filename string) ([]byte, error) {
	return s.fsUtil.ReadFile(filename)
}

func (s *directoryStorage) RootDir() string {
	return s.rootDir
}

func (s *directoryStorage) SaveUnit(u unit.Unit) error {
	if !s.IsCreated() {
		return errors.New("storage is not created yet and cannot be used")
	}
	if nil == u {
		return errors.New("cannot operate on nil unit")
	}
	return newPersistentUnit(u, *newLocation(s.RootDir(), u.ID()), s).save()
}

func (s *directoryStorage) RemoveUnit(u unit.Unit) error {
	if !s.IsCreated() {
		return errors.New("storage is not created yet and cannot be used")
	}
	if nil == u {
		return errors.New("cannot operate on nil unit")
	}
	return newPersistentUnit(u, *newLocation(s.RootDir(), u.ID()), s).remove()
}

func (s *directoryStorage) LoadUnit(id string) (unit.Unit, error) {
	if !s.IsCreated() {
		return nil, errors.New("storage is not created yet and cannot be used")
	}
	if len(id) == 0 {
		return nil, errors.New("cannot operate on nil unit")
	}
	return newPersistentUnit(nil, *newLocation(s.RootDir(), id), s).load()
}

func (s *directoryStorage) IsCreated() bool {
	_, err := s.fs.Stat(s.rootDir)

	return err == nil
}

func (s *directoryStorage) Create() error {
	return errors.Wrap(s.mkdirAll(s.rootDir, os.ModePerm), "failed to create storage")
}

func (s *directoryStorage) Remove() error {
	return errors.Wrap(s.removeDir(s.rootDir), "failed to remove storage")
}

func (s *directoryStorage) Type() string {
	return s.storageType
}

type directoryJSON struct {
	RootDir string `json:"root"`
	Type    string `json:"type"`
}

func (s *directoryStorage) MarshalJSON() ([]byte, error) {
	return json.Marshal(directoryJSON{RootDir: s.rootDir, Type: s.storageType})
}

func (s *directoryStorage) UnmarshalJSON(b []byte) error {
	var jsonData directoryJSON
	err := json.Unmarshal(b, &jsonData)

	if err != nil {
		return err
	}

	s.rootDir = jsonData.RootDir
	s.storageType = jsonData.Type

	var fs afero.Fs
	if s.storageType == TypeDirectoryInMemory {
		fs = afero.NewMemMapFs()
	} else {
		fs = afero.NewOsFs()
	}

	s.fs = fs
	s.fsUtil = &afero.Afero{Fs: fs}

	return nil
}
