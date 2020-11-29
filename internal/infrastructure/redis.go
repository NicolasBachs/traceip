package infrastructure

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	Conn *redis.Client
}

func (r *Redis) Init(addr string, pwd string) {
	r.Conn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0, // use default DB
	})

	_, err := r.Conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Could not initialize redis (", addr, ") - Error: ", err.Error())
	}
}
