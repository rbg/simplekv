package store

type mem_be struct {
	kv map[string][]byte
}

func NewMem() Store {
	return &mem_be{kv: make(map[string][]byte)}
}

func (r *mem_be) Keys() ([]string, error) {
	return nil, nil
}
func (r *mem_be) Get(key string) ([]byte, error) {
	return r.kv[key], nil
}

func (r *mem_be) Put(key string, val []byte) error {
	r.kv[key] = val
	return nil
}

func (r *mem_be) Delete(key string) error {
	delete(r.kv, key)
	return nil
}
