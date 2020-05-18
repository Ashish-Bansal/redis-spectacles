package redisscanner

import (
	"github.com/go-redis/redis"
)

var rdb *redis.Client

// GetRedisClient initialises and returns redis client from given redis URL
func GetRedisClient(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(options)
	_, err = redisClient.Ping().Result()
	return redisClient, err
}

// ScanRedisKeys scans redis database based on given pattern and sends them via channel.
func ScanRedisKeys(redisClient *redis.Client, pattern string, batchSize int64, keyReceiver chan<- string) {
	iterator := redisClient.Scan(0, pattern, batchSize).Iterator()
	for iterator.Next() {
		keyReceiver <- iterator.Val()
	}
	close(keyReceiver)
}
