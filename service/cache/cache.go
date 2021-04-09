package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var Cache = redis.NewClient(&redis.Options {
	Addr: "redis:6379",
	Password: "",
	DB: 0,
})