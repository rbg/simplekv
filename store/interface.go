package store

type Store interface {
	Keys() ([]string, error)
	Get(key string) ([]byte, error)
	Put(key string, val []byte) error
	Delete(key string) error
}
