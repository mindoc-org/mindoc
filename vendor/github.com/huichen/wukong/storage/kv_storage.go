package storage

import (
	"github.com/cznic/kv"
	"io"
)

type kvStorage struct {
	db *kv.DB
}

func openKVStorage(path string) (Storage, error) {
	options := &kv.Options{}
	db, errOpen := kv.Open(path, options)
	if errOpen != nil {
		var errCreate error
		db, errCreate = kv.Create(path, options)
		if errCreate != nil {
			return &kvStorage{db}, errCreate
		}
	}
	return &kvStorage{db}, nil
}

func (s *kvStorage) WALName() string {
	return s.db.WALName()
}

func (s *kvStorage) Set(k []byte, v []byte) error {
	return s.db.Set(k, v)
}

func (s *kvStorage) Get(k []byte) ([]byte, error) {
	return s.db.Get(nil, k)
}

func (s *kvStorage) Delete(k []byte) error {
	return s.db.Delete(k)
}

func (s *kvStorage) ForEach(fn func(k, v []byte) error) error {
	iter, err := s.db.SeekFirst()
	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}
	for {
		key, value, err := iter.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := fn(key, value); err != nil {
			return err
		}
	}
	return nil
}

func (s *kvStorage) Close() error {
	return s.db.Close()
}
