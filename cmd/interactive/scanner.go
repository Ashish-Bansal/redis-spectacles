package interactive

import (
	"fmt"

	"github.com/gdamore/tcell"
)

func displayKeysScannedMessage(screen tcell.Screen, keysScanned int) {
	message := fmt.Sprintf("Scanned %d keys...", keysScanned)
	setStringInScreen(screen, 0, 0, message, tcell.StyleDefault)
	screen.Show()
}
