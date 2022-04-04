package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

var redisCli *redis.Client

func InitRedis(addr string) error {
	redisCli = redis.NewClient(&redis.Options{Addr: addr})
	_, err := redisCli.Ping().Result()
	return err
}

func checkHistoryCache() *[]Results {
	exists, err := redisCli.Exists(historyKey).Result()
	if err != nil {
		return nil
	}
	if exists == 1 {
		value, err := redisCli.Get("ming").Result()
		if err != nil {
			return nil
		}
		result := &[]Results{}
		json.Unmarshal([]byte(value), result)
		return result
	} else {
		return nil
	}
}

func cacheData(results *[]Results) error {
	err := redisCli.Set(historyKey, results, expireTime*time.Second).Err()
	return err
}
