package store

import redis "gopkg.in/redis.v3"

type redis_be struct {
	write *redis.Client
	read  *redis.Client
}

func NewRedis(write_ep string, read_ep string) Store {

	be := &redis_be{}

	// figure out how to setup
	if len(read_ep) > 0 {
		be.read = redis.NewClient(&redis.Options{
			Addr:     read_ep,
			Password: "",
			DB:       0,
		})
		if be.read == nil {
			panic("failed to redis")
		}
		if pong, err := be.read.Ping().Result(); err != nil {
			slog.Printf("READ Ping reply: %s", pong)
		} else {
			slog.ErrPrintf("READ Ping reply error: %s", err)
			return nil
		}
	}

	if len(write_ep) == 0 {
		write_ep = "localhost:6379"
	}

	if be.write = redis.NewClient(&redis.Options{
		Addr:     write_ep,
		Password: "",
		DB:       0,
	}); be.write == nil {
		panic("failed to redis")
	}

	pong := be.write.Ping()
	if err := pong.Err(); err != nil {
		slog.ErrPrintf("WRITE Ping reply error: %s", err)
		return nil
	}

	if be.read == nil {
		be.read = be.write
	}

	return be
}

func (r *redis_be) Get(key string) ([]byte, error) {
	result := r.read.HGet("root", key)
	return result.Bytes()
}

func (r *redis_be) Put(key string, val []byte) error {
	return r.write.HSet("root", key, string(val)).Err()
}

func (r *redis_be) Delete(key string) error {
	return r.write.HDel("root", key).Err()
}

func (r *redis_be) Keys() ([]string, error) {
	return r.read.HKeys("root").Result()
}
