// Package cache
// @Author LiuYong
// @Created at 2021-02-19
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

var RedisConn *redis.Client

func Setup(ops *redis.Options) error {
	if RedisConn != nil {
		return errors.New("重复创建cache连接")
	}
	RedisConn = redis.NewClient(ops)

	return nil
}

// Set 设置一个值,expiration为0表示没有过期时间
func Set(key string, data interface{}, expiration int) error {
	ctx := context.Background()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = RedisConn.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
	return err
}

// Exists 判断关键字是否存在
func Exists(key string) bool {
	ctx := context.Background()
	// 存在返回1,不存在返回0
	exists, err := RedisConn.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return exists == 1
}

// Get 获取值,并且保存到value中
func Get(key string, value interface{}) error {
	typ := reflect.TypeOf(value)
	if typ.Kind() != reflect.Ptr {
		return errors.New("value is not a pointer")
	}
	ctx := context.Background()
	byt, err := RedisConn.Get(ctx, key).Bytes()
	if err != nil {
		return nil
	}
	return json.Unmarshal(byt, value)
}

func Delete(key string) (bool, error) {
	ctx := context.Background()
	res, err := RedisConn.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

func LikeDeletes(key string) error {
	ctx := context.Background()
	keys, err := RedisConn.Keys(ctx, "*"+key+"*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func LPush(key string, data interface{}) error {
	ctx := context.Background()
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = RedisConn.LPush(ctx, key, value).Err()
	return err
}

func LIndex(key string, index int, value interface{}) error {
	typ := reflect.TypeOf(value)
	if typ.Kind() != reflect.Ptr {
		return errors.New("value is not a pointer")
	}
	ctx := context.Background()
	bytes, err := RedisConn.LIndex(ctx, key, int64(index)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}
