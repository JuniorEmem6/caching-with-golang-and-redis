package data

import (
	"fmt"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func RedisClient() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err)
	}
	redisClient = client

}

func SetCache(id int, name string, price string, description string) bool {
	newId := fmt.Sprint(id)
	products := map[string]interface{}{
		"name":        name,
		"price":       price,
		"description": description,
	}
	resp := redisClient.HMSet(newId, products)

	return resp.Err() != nil

}

func GetCache(id int) string {
	newId := fmt.Sprint(id)

	resp := redisClient.HGetAll(newId)

	value, _ := resp.Result()
	fmt.Println(value)
	if value == nil {
		return "Cache miss"
	}

	return "Cache hit"
}
