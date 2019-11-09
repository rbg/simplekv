package store

import (
	"github.com/apex/log"
	redis "gopkg.in/redis.v3"
)

type redisBE struct {
	write *redis.Client
	read  *redis.Client
}

func NewRedis(wep string, rep string) Store {

	be := &redisBE{}

	// figure out how to setup
	if len(rep) > 0 {
		be.read = redis.NewClient(&redis.Options{
			Addr:     rep,
			Password: "",
			DB:       0,
		})
		if be.read == nil {
			panic("failed to create readonly redis")
		}
	}

	if be.write = redis.NewClient(&redis.Options{
		Addr:     wep,
		Password: "",
		DB:       0,
	}); be.write == nil {
		panic("failed to redis")
	}

	pong := be.write.Ping()
	if err := pong.Err(); err != nil {
		log.Infof("WRITE Ping reply error: %s", err)
		return nil
	}

	if be.read == nil {
		be.read = be.write
	}

	return be
}

func (r *redisBE) Get(key string) ([]byte, error) {
	result := r.read.HGet("root", key)
	return result.Bytes()
}

func (r *redisBE) Put(key string, val []byte) error {
	return r.write.HSet("root", key, string(val)).Err()
}

func (r *redisBE) Delete(key string) error {
	return r.write.HDel("root", key).Err()
}

func (r *redisBE) Keys() ([]string, error) {
	return r.read.HKeys("root").Result()
}
