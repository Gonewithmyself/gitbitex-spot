package cache

import "github.com/gitbitex/gitbitex-spot/conf"

import "github.com/go-redis/redis"

import "sync"

type cache struct {
}

func init() {
	config, _ := conf.GetConfig()
	redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
	})
}

var m sync.Map
var lock sync.Mutex

// GetClient 1
func GetClient(db int) *redis.Client {
	if v, ok := m.Load(db); ok {
		return v.(*redis.Client)
	}

	config, _ := conf.GetConfig()
	rd := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
	})
	m.Store(db, rd)
	return rd
}
