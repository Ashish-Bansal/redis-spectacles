package trie

import (
	"reflect"
	"testing"
)

func TestTrieCreation(t *testing.T) {
	node := NewNode()
	if !node.IsMutable {
		t.Error("Newly created trie node isn't mutable.")
	}

	if node.DataCount != 0 {
		t.Errorf(
			"Incorrect prefix count of the node. Expected %d, got %d",
			0,
			node.DataCount,
		)
	}
}

func TestTrieInsertion(t *testing.T) {
	node := NewNode()
	node.Insert("First")
	node.Insert("Second")

	childrenCount := len(node.Children())
	if childrenCount != 2 {
		t.Errorf(
			"Children count mismatch. Expected %d, found %d.",
			2,
			childrenCount,
		)
	}

	edges := node.Edges
	for edge := range edges {
		t.Log(edge.Prefix)
	}

	fEdge := node.GetEdge("F")
	if fEdge == nil {
		t.Errorf(
			"Edge from node not found. Expected edge with item '%s', got nil",
			"F",
		)
	}
}

func TestTrieDFS(t *testing.T) {
	node := NewNode()
	node.Insert("Bag")
	node.Insert("Bat")
	node.Insert("Boat")

	expectedPrefixes := []string{"B", "Ba", "Bag", "Bat", "Bo", "Boa", "Boat"}
	prefixes := make([]string, 0)

	node.DFS(func(item interface{}, count int) {
		prefixes = append(prefixes, item.(string))
	})

	if !reflect.DeepEqual(expectedPrefixes, prefixes) {
		t.Errorf(
			"Trie prefix mismatch. Expected %v, got %v",
			expectedPrefixes,
			prefixes,
		)
	}
}

func TestTrieBFS(t *testing.T) {
	node := NewNode()
	node.Insert("Bag")
	node.Insert("Bat")
	node.Insert("Boat")

	expectedPrefixes := []string{"B", "Ba", "Bo", "Bag", "Bat", "Boa", "Boat"}
	prefixes := make([]string, 0)

	node.BFS(func(item interface{}, count int) {
		prefixes = append(prefixes, item.(string))
	})

	if !reflect.DeepEqual(expectedPrefixes, prefixes) {
		t.Errorf(
			"Trie prefix mismatch. Expected %v, got %v",
			expectedPrefixes,
			prefixes,
		)
	}
}

func TestTrieCondensation(t *testing.T) {
	testcases := [][][]string{
		{
			{},
			{},
		},
		{
			{"Bag"},
			{"Bag"},
		},
		{
			{"Bag", "Bat", "Boat"},
			{"B", "Ba", "Boat", "Bag", "Bat"},
		},
		{
			{"The", "Bye", "Hello"},
			{"Bye", "Hello", "The"},
		},
	}

	for _, testcase := range testcases {
		stringsToBeInserted := testcase[0]
		expectedPrefixes := testcase[1]

		node := NewNode()
		for _, element := range stringsToBeInserted {
			node.Insert(element)
		}
		node.Condense()

		prefixes := make([]string, 0)

		node.BFS(func(item interface{}, count int) {
			prefixes = append(prefixes, item.(string))
		})

		if !reflect.DeepEqual(expectedPrefixes, prefixes) {
			t.Errorf(
				"Trie prefix mismatch. Expected %v, got %v",
				expectedPrefixes,
				prefixes,
			)
		}
	}
}