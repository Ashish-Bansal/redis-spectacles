package interactive

import (
	"fmt"
	"os"
	"strings"

	"container/list"

	"github.com/Ashish-Bansal/redis-spectacles/pkg/trie"
	"github.com/gdamore/tcell"
)

// ScreenRow represents information about rendering single row on screen
type ScreenRow struct {
	Message     string
	Style       tcell.Style
	PaddingLeft int
	Metadata    interface{}
}

// ScreenState represents screen state based off trie
type ScreenState struct {
	Header         []ScreenRow
	Body           []ScreenRow
	Footer         []ScreenRow
	CurrentBodyRow int
	Screen         tcell.Screen
	NodeStack      *list.List
}

func renderScreenRow(screen tcell.Screen, column int, row int, screenRow ScreenRow) {
	message := strings.Repeat(" ", screenRow.PaddingLeft) + screenRow.Message
	style := screenRow.Style
	setStringInScreen(screen, column, row, message, style)
}

func (screenState *ScreenState) render() {
	screen := screenState.Screen
	screen.Clear()

	for index, screenRow := range screenState.Header {
		row := index
		renderScreenRow(screen, 0, row, screenRow)
		setRowBackground(screen, row, screenRow.Style)
	}

	headerLength := len(screenState.Header)
	for index, screenRow := range screenState.Body {
		row := index + headerLength
		renderScreenRow(screen, 0, row, screenRow)
	}

	if len(screenState.Body) != 0 {
		setRowBackground(screen, len(screenState.Header), highlighedStyle)
	}

	_, height := screen.Size()
	footerLength := len(screenState.Footer)
	for index, screenRow := range screenState.Footer {
		row := height + index - footerLength
		renderScreenRow(screen, 0, row, screenRow)
		setRowBackground(screen, row, screenRow.Style)
	}

	screenState.Screen.Show()
}

func getHeader(node *trie.Node) []ScreenRow {
	header := []ScreenRow{
		{
			Message: "redis-spectacles ~ Use the arrow keys to navigate.",
			Style:   highlighedStyle,
		},
		{
			Message: strings.Repeat("-", 500),
			Style:   normalStyle,
		},
	}
	return header
}

func getFooter(node *trie.Node) []ScreenRow {
	footer := []ScreenRow{
		{
			Message:     fmt.Sprintf("Total key count : %d", node.Count()),
			Style:       highlighedStyle,
			PaddingLeft: 1,
		},
	}
	return footer
}

func updateTrieNodeInScreenState(screenState *ScreenState, rootNode *trie.Node) {
	body := make([]ScreenRow, 0)
	for _, edge := range rootNode.GetEdges() {
		node := rootNode.Edges[edge]
		message := fmt.Sprintf("%s - %d", edge.Prefix, node.Count())
		row := ScreenRow{Message: message, Style: normalStyle, PaddingLeft: 5, Metadata: node}
		body = append(body, row)
	}
	screenState.Body = body
	screenState.CurrentBodyRow = 0
	screenState.render()
}

func popNodeFromStack(screenState *ScreenState) {
	nodeStack := screenState.NodeStack
	if nodeStack.Len() < 2 {
		return
	}

	topNodeElement := nodeStack.Back()
	nodeStack.Remove(topNodeElement)
	topNodeElement = nodeStack.Back()

	node := topNodeElement.Value.(*trie.Node)
	updateTrieNodeInScreenState(screenState, node)
}

func pushNodeIntoStack(screenState *ScreenState, node *trie.Node) {
	nodeStack := screenState.NodeStack
	nodeStack.PushBack(node)

	updateTrieNodeInScreenState(screenState, node)
}

func initScreenState(screen tcell.Screen, node *trie.Node) *ScreenState {
	header := getHeader(node)
	footer := getFooter(node)
	screenState := ScreenState{Screen: screen, Header: header, Footer: footer, NodeStack: list.New()}
	pushNodeIntoStack(&screenState, node)
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

func handleKeyLeft(screenState *ScreenState) {
	popNodeFromStack(screenState)
}

func handleKeyRight(screenState *ScreenState) {
	screenRow := screenState.Body[screenState.CurrentBodyRow]
	node := screenRow.Metadata.(*trie.Node)
	if len(node.Edges) == 0 {
		return
	}

	pushNodeIntoStack(screenState, node)
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
	case tcell.KeyLeft:
		handleKeyLeft(screenState)
	case tcell.KeyRight:
		handleKeyRight(screenState)
	}
}

func handleResize(screenState *ScreenState, event *tcell.EventResize) {
	screenState.Screen.Clear()
	screenState.render()
	// ToDo: Handle min terminal size
}
