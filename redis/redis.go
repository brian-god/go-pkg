package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func ConfigDatabase(opt Option) (err error) {
	option = &opt
	DB = redis.NewClient(&redis.Options{
		Addr:     option.Addr,
		Password: option.Password,
		DB:       option.DB,
	})
	_ = DB

	pool := goredis.NewPool(DB)
	RS = redsync.New(pool)
	_ = RS

	return nil
}

// MutexWithUnlock 分布式锁，并发控制
func MutexWithUnlock(name string, options ...redsync.Option) (UnlockFunc, error) {
	mutex := RS.NewMutex(name, options...)
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
func SimpleMutexWithUnlock(name string) (UnlockFunc, error) {
	mutex := RS.NewMutex(name)
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
	return err != nil && err != redis.Nil
}
