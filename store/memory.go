package store

import "fmt"

type mem_be struct {
	kv map[string][]byte
}

func NewMem() Store {
	be_mem := &mem_be{kv: make(map[string][]byte)}
	slog.Printf("NewStore: %p", be_mem)
	return be_mem
}

func (r *mem_be) Keys() ([]string, error) {
	var keys []string

	slog.Printf("Keys: %p", r)
	for k := range r.kv {
		keys = append(keys, k)
	}
	return keys, nil
}

func (r *mem_be) Get(key string) ([]byte, error) {
	if val, ok := r.kv[key]; !ok {
		return nil, fmt.Errorf("No value for key")
	} else {
		slog.Printf("Get: %p %s %s", r, key, string(val))
		return val, nil
	}
}

func (r *mem_be) Put(key string, val []byte) error {
	slog.Printf("Put: %p %s %s", r, key, string(val))
	r.kv[key] = val
	return nil
}

func (r *mem_be) Delete(key string) error {
	slog.Printf("Delete: %p %s", r, key)
	delete(r.kv, key)
	return nil
}
