package interactive

import (
	"log"

	"github.com/urfave/cli/v2"

	"github.com/Ashish-Bansal/redis-spectacles/internal/consts"
	"github.com/Ashish-Bansal/redis-spectacles/internal/redisscanner"
	"github.com/Ashish-Bansal/redis-spectacles/pkg/trie"
)

// ExecuteInteractive runs interactive version of redis analyzer
func ExecuteInteractive(c *cli.Context) {
	redisURL := c.String(consts.RedisURLArgName)
	scanBatchSize := c.Int64(consts.ScanBatchSizeArgName)
	scanPattern := c.String(consts.ScanPattern)

	client, err := redisscanner.GetRedisClient(redisURL)
	if err != nil {
		log.Fatal(err)
	}

	keyReceiver := make(chan string)
	go redisscanner.ScanRedisKeys(client, scanPattern, scanBatchSize, keyReceiver)

	screen := initScreen()
	node := trie.NewNode()
	keysScanned := 0
	for key := range keyReceiver {
		node.Insert(key)
		screen.Clear()
		keysScanned++
		displayKeysScannedMessage(screen, keysScanned)
	}
	node.Condense()
	screen.Clear()

	screenState := initScreenState(screen, node)
	startEventLoop(screenState)
}
