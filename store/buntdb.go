package store

import (
	"github.com/apex/log"
	"github.com/tidwall/buntdb"
)

type buntdbBE struct {
	db *buntdb.DB
}

func NewBuntDB() Store {
	db, _ := buntdb.Open(":memory:")
	beBuntDB := &buntdbBE{db: db}
	log.Debugf("NewStore: %p", beBuntDB)
	return beBuntDB
}

func (r *buntdbBE) Keys() (keys []string, err error) {
	var k []string

	err = r.db.View(func(tx *buntdb.Tx) error {
		tx.AscendKeys("*", func(key, val string) bool {
			k = append(k, key)
			return true
		})
		if err != nil {
			return err
		}
		return nil
	})

	return k, nil
}

func (r *buntdbBE) Get(key string) (val []byte, err error) {
	var dbVal string
	log.Debugf("Get: %p %s ", r, key)
	err = r.db.View(func(tx *buntdb.Tx) error {
		dbVal, err = tx.Get(key)
		if err != nil {
			return err
		}
		return nil
	})
	b := []byte(dbVal)
	return b, err
}

func (r *buntdbBE) Put(key string, val []byte) error {
	log.Debugf("Put: %p %s %s", r, key, string(val))
	err := r.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(val), nil)
		return err
	})
	return err
}

func (r *buntdbBE) Delete(key string) error {
	log.Debugf("Delete: %p %s", r, key)
	err := r.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
	return err
}
