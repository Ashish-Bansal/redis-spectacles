package trie

import (
	"errors"
	"sort"

	"github.com/Ashish-Bansal/redis-spectacles/pkg/addable"
	"github.com/Ashish-Bansal/redis-spectacles/pkg/comparable"
	"github.com/Ashish-Bansal/redis-spectacles/pkg/iterator"
)

// Edge is the connection between two nodes
type Edge struct {
	PrefixCount int
	Prefix      interface{}
}

// Node implementing NodeInterface
type Node struct {
	Edges     map[*Edge]*Node
	IsMutable bool
	DataCount int
}

// WalkCallback is a callback for the bfs/dfs functions
type WalkCallback func(interface{}, int)

// NewNode creates new trie node
func NewNode() *Node {
	return &Node{IsMutable: true, Edges: make(map[*Edge]*Node)}
}

// Children returns all direct children for a trie node
func (node *Node) Children() []*Node {
	edges := node.Edges
	children := []*Node{}
	for _, value := range edges {
		children = append(children, value)
	}
	return children
}

// Count returns sum of items passing through this node i.e. prefix and data counts
func (node *Node) Count() int {
	edges := node.Edges
	count := node.DataCount
	for edge := range edges {
		count += edge.PrefixCount
	}
	return count
}

// GetEdges returns edges of the node in the sorted order of prefixes
func (node *Node) GetEdges() []*Edge {
	edges := node.Edges
	orderedEdges := make([]*Edge, 0)
	for edge := range edges {
		orderedEdges = append(orderedEdges, edge)
	}
	sort.Slice(orderedEdges, func(i int, j int) bool {
		a := orderedEdges[i]
		b := orderedEdges[j]
		result, err := comparable.LessThan(a.Prefix, b.Prefix)
		if err != nil {
			panic(err)
		}
		return result
	})
	return orderedEdges
}

// GetEdge returns edge if it exists from current node, otherwise returns nil
func (node *Node) GetEdge(item interface{}) *Edge {
	edges := node.Edges
	for edge := range edges {
		if edge.Prefix == item {
			return edge
		}
	}
	return nil
}

// Insert adds new element into trie
func (node *Node) Insert(prefix interface{}) error {
	if !node.IsMutable {
		return errors.New("Trying to run insert on non-mutable tree instance")
	}

	iterator, err := iterator.NewIterator(prefix)
	if err != nil {
		return errors.New("Don't know how to iterate passed object")
	}

	currentNode := node
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			panic(err)
		}

		edge := currentNode.GetEdge(item)
		if edge == nil {
			edge = &Edge{Prefix: item}
			currentNode.Edges[edge] = NewNode()
		}

		edge.PrefixCount++
		currentNode = currentNode.Edges[edge]
	}
	currentNode.DataCount++
	return nil
}

// Condense marks the trie as un-mutable and in case parent has single child, it merges itself with parent node.
// Prefix interface must support Addition operation, otherwise it will cause panic.
func (node *Node) Condense() error {
	if !node.IsMutable {
		return errors.New("Trying to run Condense on non-mutable tree instance")
	}

	for childEdge, child := range node.Edges {
		child.Condense()
		childEdges := child.Edges
		if len(childEdges) == 1 {
			for grandChildEdge, grandChildNode := range childEdges {
				newKey, err := addable.Add(childEdge.Prefix, grandChildEdge.Prefix)
				newEdge := &Edge{Prefix: newKey, PrefixCount: 1}
				if err != nil {
					panic(err)
				}
				delete(node.Edges, childEdge)
				node.Edges[newEdge] = grandChildNode
			}
		}
	}

	node.IsMutable = false
	return nil
}

func (node *Node) dfs(callback WalkCallback, prefix interface{}) {
	for _, edge := range node.GetEdges() {
		newPrefix, _ := addable.Add(prefix, edge.Prefix)
		callback(newPrefix, edge.PrefixCount)
		childNode := node.Edges[edge]
		childNode.dfs(callback, newPrefix)
	}
}

// DFS performs breadth first search on the trie and calls callback with each node element and prefix till now
func (node *Node) DFS(callback WalkCallback) {
	node.dfs(callback, nil)
}

// BFS performs breadth first search on the trie and calls callback with each node element
func (node *Node) BFS(callback WalkCallback) {
	type Pair struct {
		node   *Node
		prefix interface{}
	}

	queue := make([]Pair, 0)
	queue = append(queue, Pair{node: node, prefix: ""})

	for len(queue) != 0 {
		pair := queue[0]
		queue = queue[1:]

		currentNode := pair.node
		prefix := pair.prefix

		for _, edge := range currentNode.GetEdges() {
			newPrefix, _ := addable.Add(prefix, edge.Prefix)
			callback(newPrefix, edge.PrefixCount)

			childNode := currentNode.Edges[edge]
			queue = append(queue, Pair{node: childNode, prefix: newPrefix})
		}
	}
}
