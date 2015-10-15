package store

import "github.com/stackengine/selog"

var slog = selog.Register("store", 0)

type Store interface {
	Keys() ([]string, error)
	Get(key string) ([]byte, error)
	Put(key string, val []byte) error
	Delete(key string) error
}
