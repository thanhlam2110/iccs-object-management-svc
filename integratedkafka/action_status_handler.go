package integratedkafka

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/cloud-platform/vnpt-sso-usermgnt/config"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

//var rdb *redis.Client
//redis
var ctxRedis = context.Background()

func ConnectRedis() (rdb *redis.Client, err error) {
	config.ReadConfig()
	url := viper.GetString(`redis.url`)
	rdb = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "ReDis0rimx", // no password set
		DB:       0,            // use default DB
	})
	_, err = rdb.Ping(ctxRedis).Result()
	return rdb, err
}
func KafkaCreateStatus(key, status string) error {
	rdb, err := ConnectRedis()
	if err != nil {
		return err
	}
	err = rdb.SetNX(ctxRedis, key, status, 7200*time.Second).Err()
	return nil
}
func KafkaReadStatus(key string) (status string, err error) {
	rdb, err := ConnectRedis()
	if err != nil {
		return "", err
	}
	val, err := rdb.Get(ctxRedis, key).Result()
	if err == redis.Nil {
		fmt.Println("key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println(key, val)
	}
	return val, nil
}
func KafkaDeleteStatus(key string) error {
	rdb, err := ConnectRedis()
	if err != nil {
		return err
	}
	err = rdb.Del(ctxRedis, key).Err()
	if err != nil {
		return err
	}
	return nil
}
