package store

type redis_be struct {
}

func NewRedis() Store {
	return &redis_be{}
}

func (r *redis_be) Get(key string, val []byte) error {
	return nil
}

func (r *redis_be) Put(key string, val []byte) error {
	return nil
}

func (r *redis_be) Delete(key string) error {
	return nil
}
