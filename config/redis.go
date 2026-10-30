package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/bsm/redislock"
	"github.com/linn221/expenses/utils"
	"github.com/redis/go-redis/v9"
)

var (
	_ utils.CacheService = (*RedisService)(nil)
	_ utils.Counter      = (*RedisService)(nil)
)

type RedisService struct {
	ctx    context.Context
	rdb    *redis.Client
	locker *redislock.Client
}

func NewRedisService(ctx context.Context, rdb *redis.Client) *RedisService {
	locker := redislock.New(rdb)
	return &RedisService{
		ctx:    ctx,
		rdb:    rdb,
		locker: locker,
	}
}

func ConnectRedis(ctx context.Context) *RedisService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("Fail to Connect Redis")
	}
	return NewRedisService(ctx, rdb)
}

func (rc *RedisService) GetClient() *redis.Client {
	return rc.rdb
}

func (rc *RedisService) GetObject(key string, dest any) (bool, error) {
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	err = json.Unmarshal([]byte(val), &dest)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rc *RedisService) GetValue(key string) (string, bool, error) {
	val, err := rc.rdb.Get(rc.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", false, nil
		}
		return "", false, err
	}
	return val, true, nil
}

func (rc *RedisService) SetObject(key string, obj any, exp time.Duration) error {
	objInByte, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if err = rc.rdb.Set(rc.ctx, key, objInByte, exp).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisService) SetValue(key string, value string, exp time.Duration) error {
	return rc.rdb.Set(rc.ctx, key, value, exp).Err()
}

func (rc *RedisService) RemoveKey(ck string) error {
	_, err := rc.rdb.Del(rc.ctx, ck).Result()
	return err
}

func (rc *RedisService) RemoveKeyWithCount(ck string) (int64, error) {
	return rc.rdb.Del(rc.ctx, ck).Result()
}

func (rc *RedisService) RemoveKeysWithCount(cks []string) (int64, error) {
	return rc.rdb.Del(rc.ctx, cks...).Result()
}

func (rc *RedisService) AddSet(setKey string, member string) error {
	if err := rc.rdb.SAdd(rc.ctx, setKey, member).Err(); err != nil {
		return err
	}
	return nil
}

func (rc *RedisService) GetSetMembers(setKey string) ([]string, error) {
	return rc.rdb.SMembers(rc.ctx, setKey).Result()
}

func (rc *RedisService) RemoveSetMember(setKey string, member string) error {
	return rc.rdb.SRem(rc.ctx, setKey, member).Err()
}

func (rc *RedisService) RemoveKeys(keys ...string) error {
	k := make([]string, 0, len(keys))
	for _, ck := range keys {
		k = append(k, ck)
	}
	_, err := rc.rdb.Del(rc.ctx, k...).Result()
	return err
}

func (rc *RedisService) SetH(key string, v map[string]any, exp time.Duration) error {
	_, err := rc.rdb.TxPipelined(rc.ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(rc.ctx, key, v)
		pipe.Expire(rc.ctx, key, exp)
		return nil
	})
	return err
}

func (rc *RedisService) GetH(key string, field string) (string, error) {
	return rc.rdb.HGet(rc.ctx, key, field).Result()
}

func (rc *RedisService) Exists(key string) (int64, error) {
	return rc.rdb.Exists(rc.ctx, key).Result()
}

func (rc *RedisService) Incr(key string) (int64, error) {
	return rc.rdb.Incr(rc.ctx, key).Result()
}

func (rc *RedisService) Lock(lock string) (func() error, error) {
	l, err := rc.locker.Obtain(rc.ctx, lock, Cfg.LockTTL(), nil)
	if err != nil {
		return nil, err
	}
	return func() error {
		return l.Release(rc.ctx)
	}, nil
}
