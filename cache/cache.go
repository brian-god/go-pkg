package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	IsNotFound      = errors.New("cache is not found")
	MarshalFailed   = errors.New("cache marshal failed")
	UnmarshalFailed = errors.New("cache unmarshal failed")
)

// GetCacheData 获取缓存数据
func GetCacheData[T any](ctx context.Context, rc *rockscache.Client, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	var t T
	var write bool
	v, err := rc.Fetch2(ctx, key, expire, func() (s string, err error) {
		t, err = fn(ctx)
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		write = true
		return string(bs), nil
	})
	if err != nil {
		return t, err
	}
	if write {
		return t, nil
	}
	if v == "" {
		return t, IsNotFound
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		return t, err
	}
	return t, nil
}

// BatchGetCache 批量获取缓存
func BatchGetCache[T any](ctx context.Context, rcClient *rockscache.Client, keys []string, expire time.Duration, keyIndexFn func(t T, keys []string) (int, error), fn func(ctx context.Context) ([]T, error)) ([]T, error) {
	batchMap, err := rcClient.FetchBatch2(ctx, keys, expire, func(idxs []int) (m map[int]string, err error) {
		values := make(map[int]string)
		tArrays, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range tArrays {
			index, err := keyIndexFn(v, keys)
			if err != nil {
				continue
			}
			bs, err := json.Marshal(v)
			if err != nil {
				return nil, MarshalFailed
			}
			values[index] = string(bs)
		}
		return values, nil
	})
	if err != nil {
		return nil, err
	}
	var tArrays []T
	for _, v := range batchMap {
		if v != "" {
			var t T
			err = json.Unmarshal([]byte(v), &t)
			if err != nil {
				return nil, UnmarshalFailed
			}
			tArrays = append(tArrays, t)
		}
	}
	return tArrays, nil
}

// DeletePrefixData 根据前缀删除
func DeletePrefixData(ctx context.Context, rd *redis.Client, rc *rockscache.Client, prefix string) error {
	result, err := rd.Keys(ctx, prefix).Result()
	if err != nil {
		return err
	}
	if len(result) > 0 {
		return rc.TagAsDeletedBatch2(ctx, result)
	}
	return nil
}
