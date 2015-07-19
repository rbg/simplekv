package store

type Store interface {
	Get(key string, val []byte) error
	Put(key string, val []byte) error
	Delete(key string) error
}
