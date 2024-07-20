package redis

import (
	"context"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	rs     *redsync.Redsync
}

func NewRedisClient(opt Option) (*RedisClient, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Password: opt.Password,
		DB:       opt.DB,
	})
	err := db.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	pool := goredis.NewPool(db)
	rs := redsync.New(pool)
	return &RedisClient{
		client: db,
		rs:     rs,
	}, nil
}

// MutexWithUnlock 分布式锁，并发控制
func (rc *RedisClient) MutexWithUnlock(name string, options ...redsync.Option) (UnlockFunc, error) {
	mutex := rc.rs.NewMutex(name, options...)
	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	unlock := func() error {
		_, err := mutex.Unlock()
		if err != nil {
			return err
		}
		return nil
	}

	return unlock, nil
}

// SimpleMutexWithUnlock 分布式锁，并发控制
func (rc *RedisClient) SimpleMutexWithUnlock(name string) (UnlockFunc, error) {
	mutex := rc.rs.NewMutex(name)
	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	unlock := func() error {
		_, err := mutex.Unlock()
		if err != nil {
			return err
		}
		return nil
	}

	return unlock, nil
}

// IfErrorNotNil 是否为非空错误
func IfErrorNotNil(err error) bool {
	return err != nil && !errors.Is(err, redis.Nil)
}
