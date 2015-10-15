package store

import "gopkg.in/redis.v3"

type redis_be struct {
	clnt *redis.Client
}

func NewRedis() Store {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:7008",
		Password: "",
		DB:       0,
	})

	if client == nil {
		panic("failed to redis")
	}
	pong, err := client.Ping().Result()
	slog.Println(pong, err)

	return &redis_be{
		clnt: client,
	}
}

func (r *redis_be) Get(key string) ([]byte, error) {
	result := r.clnt.HGet("root", key)
	return result.Bytes()
}

func (r *redis_be) Put(key string, val []byte) error {
	return r.clnt.HSet("root", key, string(val)).Err()
}

func (r *redis_be) Delete(key string) error {
	return r.clnt.HDel("root", key).Err()
}

func (r *redis_be) Keys() ([]string, error) {
	return r.clnt.HKeys("root").Result()
}
