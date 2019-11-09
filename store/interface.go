package store

// Store represents what each store implementaion should
// at least provide, it is an interface after all..
type Store interface {
	Keys() ([]string, error)
	Get(key string) ([]byte, error)
	Put(key string, val []byte) error
	Delete(key string) error
}
