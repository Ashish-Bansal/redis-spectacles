package interactive

import (
	"fmt"
	"os"

	"github.com/Ashish-Bansal/redis-spectacles/pkg/trie"
	"github.com/gdamore/tcell"
)

// ScreenState represents screen state based off trie
type ScreenState struct {
	Header         []string
	Body           []string
	Footer         []string
	CurrentBodyRow int
	Screen         tcell.Screen
}

func (screenState *ScreenState) render() {
	style := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	for index, message := range screenState.Header {
		row := index
		setStringInScreen(screenState.Screen, 0, row, message, style)
		setRowBackground(screenState.Screen, row, style)
	}

	style = tcell.StyleDefault
	headerLength := len(screenState.Header)
	for index, message := range screenState.Body {
		row := index + headerLength
		setStringInScreen(screenState.Screen, 0, row, message, style)
	}

	_, height := screenState.Screen.Size()
	style = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	footerLength := len(screenState.Footer)
	for index, message := range screenState.Footer {
		row := height + index - footerLength
		setStringInScreen(screenState.Screen, 0, row, message, style)
		setRowBackground(screenState.Screen, row, style)
	}

	screenState.Screen.Show()
}

func getHeader(node *trie.Node) []string {
	header := make([]string, 0)
	header = append(header, "redis-spectacles ~ Use the arrow keys to navigate.")
	return header
}

func getFooter(node *trie.Node) []string {
	footer := make([]string, 0)
	footer = append(footer, fmt.Sprintf("Total key count : %d", node.Count()))
	return footer
}

func updateTrieNodeInScreenState(screenState *ScreenState, node *trie.Node) {
	body := make([]string, 0)
	for edge, node := range node.Edges {
		message := fmt.Sprintf("%s - %d", edge.Prefix, node.Count())
		body = append(body, message)
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

func handleKeyEvent(screenState *ScreenState, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyCtrlC:
		screenState.Screen.Fini()
		os.Exit(0)
	}
}

func handleResize(screenState *ScreenState, event *tcell.EventResize) {
	screenState.Screen.Clear()
	screenState.render()
	// ToDo: Handle min terminal size
}
