package config

import (
	"mango-admin/pkg/storage"
	"mango-admin/pkg/storage/cache"
)

type Cache struct {
	Redis  *RedisConnectOptions
	Memory interface{}
}

// CacheConfig cache配置
var CacheConfig = new(Cache)

// Setup 构造cache 顺序 redis > 其他 > memory
func (e Cache) Setup() (storage.AdapterCache, error) {
	if e.Redis != nil {
		options, err := e.Redis.GetRedisOptions()
		if err != nil {
			return nil, err
		}
		r, err := cache.NewRedis(GetRedisClient(), options)
		if err != nil {
			return nil, err
		}
		if _redis == nil {
			_redis = r.GetClient()
		}
		return r, nil
	}
	return cache.NewMemory(), nil
}
