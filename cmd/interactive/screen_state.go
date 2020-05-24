package interactive

import (
	"fmt"
	"os"

	"github.com/Ashish-Bansal/redis-spectacles/pkg/trie"
	"github.com/gdamore/tcell"
)

// ScreenRow represents information about rendering single row on screen
type ScreenRow struct {
	Message string
	Style   tcell.Style
}

// ScreenState represents screen state based off trie
type ScreenState struct {
	Header         []ScreenRow
	Body           []ScreenRow
	Footer         []ScreenRow
	CurrentBodyRow int
	Screen         tcell.Screen
}

func (screenState *ScreenState) render() {
	for index, screenRow := range screenState.Header {
		row := index
		setStringInScreen(screenState.Screen, 0, row, screenRow.Message, screenRow.Style)
		setRowBackground(screenState.Screen, row, highlighedStyle)
	}

	headerLength := len(screenState.Header)
	for index, screenRow := range screenState.Body {
		row := index + headerLength
		setStringInScreen(screenState.Screen, 0, row, screenRow.Message, screenRow.Style)
	}

	if len(screenState.Body) != 0 {
		setRowBackground(screenState.Screen, len(screenState.Header), highlighedStyle)
	}

	_, height := screenState.Screen.Size()
	footerLength := len(screenState.Footer)
	for index, screenRow := range screenState.Footer {
		row := height + index - footerLength
		setStringInScreen(screenState.Screen, 0, row, screenRow.Message, screenRow.Style)
		setRowBackground(screenState.Screen, row, highlighedStyle)
	}

	screenState.Screen.Show()
}

func getHeader(node *trie.Node) []ScreenRow {
	header := []ScreenRow{
		{
			Message: "redis-spectacles ~ Use the arrow keys to navigate.",
			Style:   highlighedStyle,
		},
	}
	return header
}

func getFooter(node *trie.Node) []ScreenRow {
	footer := []ScreenRow{
		{
			Message: fmt.Sprintf("Total key count : %d", node.Count()),
			Style:   highlighedStyle,
		},
	}
	return footer
}

func updateTrieNodeInScreenState(screenState *ScreenState, node *trie.Node) {
	body := make([]ScreenRow, 0)
	for edge, node := range node.Edges {
		message := fmt.Sprintf("%s - %d", edge.Prefix, node.Count())
		row := ScreenRow{Message: message, Style: normalStyle}
		body = append(body, row)
	}
	screenState.Body = body
	screenState.render()
}

func initScreenState(screen tcell.Screen, node *trie.Node) *ScreenState {
	header := getHeader(node)
	footer := getFooter(node)
	screenState := ScreenState{Screen: screen, Header: header, Footer: footer}
	updateTrieNodeInScreenState(&screenState, node)
	return &screenState
}

func startEventLoop(screenState *ScreenState) {
	screen := screenState.Screen
	for {
		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			handleKeyEvent(screenState, event)
		case *tcell.EventResize:
			handleResize(screenState, event)
		}
	}
}

func handleKeyDown(screenState *ScreenState) {
	screen := screenState.Screen
	currentBodyRow := screenState.CurrentBodyRow
	newBodyRow := currentBodyRow + 1

	maxBodySize := len(screenState.Body)
	if newBodyRow == maxBodySize {
		return
	}

	currentRow := currentBodyRow + len(screenState.Header)
	newRow := newBodyRow + len(screenState.Header)

	setRowBackground(screen, currentRow, normalStyle)
	setRowBackground(screen, newRow, highlighedStyle)
	screen.Show()
	screenState.CurrentBodyRow++
}

func handleKeyUp(screenState *ScreenState) {
	screen := screenState.Screen
	currentBodyRow := screenState.CurrentBodyRow
	newBodyRow := currentBodyRow - 1
	if newBodyRow < 0 {
		return
	}

	currentRow := currentBodyRow + len(screenState.Header)
	newRow := newBodyRow + len(screenState.Header)

	setRowBackground(screen, currentRow, normalStyle)
	setRowBackground(screen, newRow, highlighedStyle)
	screen.Show()
	screenState.CurrentBodyRow--
}

func handleKeyEvent(screenState *ScreenState, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyRune:
		if event.Rune() == 'q' {
			screenState.Screen.Fini()
			os.Exit(0)
		}
	case tcell.KeyCtrlC:
		screenState.Screen.Fini()
		os.Exit(1)
	case tcell.KeyUp:
		handleKeyUp(screenState)
	case tcell.KeyDown:
		handleKeyDown(screenState)
	}
}

func handleResize(screenState *ScreenState, event *tcell.EventResize) {
	screenState.Screen.Clear()
	screenState.render()
	// ToDo: Handle min terminal size
}
