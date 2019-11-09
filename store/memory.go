package store

import (
	"fmt"

	"github.com/apex/log"
)

type memBE struct {
	kv map[string][]byte
}

func NewMem() Store {
	be := &memBE{kv: make(map[string][]byte)}
	log.Debugf("NewStore: %+#v", be)
	return be
}

func (r *memBE) Keys() ([]string, error) {
	var keys []string

	log.Debugf("Keys: %+#v", r)
	for k := range r.kv {
		keys = append(keys, k)
	}
	return keys, nil
}

func (r *memBE) Get(key string) ([]byte, error) {
	if val, ok := r.kv[key]; !ok {
		return nil, fmt.Errorf("No value for key")
	} else {
		log.Debugf("Get: %p %s %s", r, key, string(val))
		return val, nil
	}
}

func (r *memBE) Put(key string, val []byte) error {
	log.Debugf("Put: %p %s %s", r, key, string(val))
	r.kv[key] = val
	return nil
}

func (r *memBE) Delete(key string) error {
	log.Debugf("Delete: %p %s", r, key)
	delete(r.kv, key)
	return nil
}
