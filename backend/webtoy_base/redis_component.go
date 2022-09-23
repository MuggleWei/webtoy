package webtoy_base

import (
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

type RedisComponent struct {
	Client *redis.Client
}

var (
	instRedisComponent *RedisComponent
	onceRedisComponent sync.Once
)

func GetRedisComponent() *RedisComponent {
	if instRedisComponent == nil {
		onceRedisComponent.Do(func() {
			instRedisComponent = &RedisComponent{
				Client: nil,
			}
		})
	}
	return instRedisComponent
}

func (this *RedisComponent) Init(host string, port uint, passwd string, db int) error {
	redisAddr := host + ":" + strconv.Itoa(int(port))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: passwd,
		DB:       db,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}

	this.Client = redisClient

	return nil
}
