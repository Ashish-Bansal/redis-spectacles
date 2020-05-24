package interactive

import (
	"fmt"
	"os"

	"github.com/Ashish-Bansal/redis-spectacles/internal/utils"
	"github.com/gdamore/tcell"
)

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

func setStringInScreen(screen tcell.Screen, columnStart int, row int, message string, style tcell.Style) {
	width, _ := screen.Size()
	length := utils.Min(len(message), width)
	for index := 0; index < length; index++ {
		character := rune(message[index])
		screen.SetContent(columnStart, row, character, []rune(""), style)
		columnStart++
	}
	screen.Show()
}

func setRowBackground(screen tcell.Screen, row int, style tcell.Style) {
	width, _ := screen.Size()
	for column := 0; column < width; column++ {
		character, comb, _, _ := screen.GetContent(column, row)
		screen.SetContent(column, row, character, comb, style)
	}
}
