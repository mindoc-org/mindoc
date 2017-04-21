package storage

import (
	"fmt"
	"os"
)

const DEFAULT_STORAGE_ENGINE = "bolt"

var supportedStorage = map[string]func(path string) (Storage, error){
	"kv":   openKVStorage,
	"bolt": openBoltStorage,
}

func RegisterStorageEngine(name string, fn func(path string) (Storage, error)) {
	supportedStorage[name] = fn
}

type Storage interface {
	Set(k, v []byte) error
	Get(k []byte) ([]byte, error)
	Delete(k []byte) error
	ForEach(fn func(k, v []byte) error) error
	Close() error
	WALName() string
}

func OpenStorage(path string) (Storage, error) {
	wse := os.Getenv("WUKONG_STORAGE_ENGINE")
	if wse == "" {
		wse = DEFAULT_STORAGE_ENGINE
	}
	if fn, has := supportedStorage[wse]; has {
		return fn(path)
	}
	return nil, fmt.Errorf("unsupported storage engine %v", wse)
}
