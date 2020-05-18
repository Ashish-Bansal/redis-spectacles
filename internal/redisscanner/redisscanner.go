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

// ScanRedisKeys scans redis database based on given pattern and provides progress updates using channel.
func ScanRedisKeys(redisClient *redis.Client, pattern string, batchSize int64, keysScannedProgress chan<- int) []string {
	allKeys := make([]string, 0)
	totalKeys := 0

	iterator := redisClient.Scan(0, pattern, batchSize).Iterator()
	for iterator.Next() {
		allKeys = append(allKeys, iterator.Val())
		totalKeys++
		keysScannedProgress <- totalKeys
	}
	return allKeys
}
