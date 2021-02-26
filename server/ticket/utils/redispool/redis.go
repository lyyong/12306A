// @Author: KongLingWen
// @Created at 2021/2/25
// @Modified at 2021/2/25

package redispool

import (
	"github.com/gomodule/redigo/redis"
	"ticket/utils/setting"
)

var RedisPool *redis.Pool

func init(){
	RedisPool = newRedisPool()
}

func newRedisPool() *redis.Pool {
	return &redis.Pool {
		MaxIdle: setting.Redis.MaxIdle,
		IdleTimeout: setting.Redis.IdleTimeout,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", setting.Redis.Host)
			if err != nil {
				return nil, err
			}
			return conn, err
		},
	}
}
