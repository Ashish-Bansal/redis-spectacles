package noninteractive

import (
	"fmt"
	"log"

	"github.com/Ashish-Bansal/redis-spectacles/internal/consts"
	"github.com/Ashish-Bansal/redis-spectacles/internal/redisscanner"
	"github.com/Ashish-Bansal/redis-spectacles/pkg/trie"
	"github.com/urfave/cli/v2"
)

func ExecuteNonInteractive(c *cli.Context) {
	redisURL := c.String(consts.RedisURLArgName)
	scanBatchSize := c.Int64(consts.ScanBatchSizeArgName)
	scanPattern := c.String(consts.ScanPattern)

	client, err := redisscanner.GetRedisClient(redisURL)
	if err != nil {
		log.Fatal(err)
	}

	keyReceiver := make(chan string, 100)
	go redisscanner.ScanRedisKeys(client, scanPattern, scanBatchSize, keyReceiver)

	node := trie.NewNode()
	for key := range keyReceiver {
		node.Insert(key)
	}
	node.Condense()

	prefixes := make([]string, 0)
	node.DFS(func(item interface{}, count int) {
		prefixes = append(prefixes, item.(string))
	})
	fmt.Println(prefixes)
}
