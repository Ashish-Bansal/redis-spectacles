package interactive

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/urfave/cli/v2"

	"github.com/Ashish-Bansal/redis-spectacles/internal/consts"
	"github.com/Ashish-Bansal/redis-spectacles/internal/redisscanner"
	"github.com/Ashish-Bansal/redis-spectacles/pkg/trie"
)

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

	screen.Fini()
}

func displayKeysScannedMessage(screen tcell.Screen, keysScanned int) {
	message := fmt.Sprintf("Scanned %d keys...", keysScanned)

	column := 0
	row := 0

	for _, c := range message {
		screen.SetContent(column, row, c, []rune(""), tcell.StyleDefault)
		column++
	}
	screen.Show()
}

func initScreen() tcell.Screen {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
	screen.Clear()
	return screen
}
