package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

var (
	DB     *redis.Client
	RS     *redsync.Redsync
	option *Option
)

type UnlockFunc func() error

type Option struct {
	Addr     string
	Password string
	DB       int
}
