package trie

import (
	"bufio"
	"os"
	"strings"
)

// Node contains a Boolean value to determine whether
// or not it terminates a word, and a map that associates runes with Nodes.
// If Terminal is false, Children cannot be empty.
type Node struct {
	Terminal bool
	Children map[rune]*Node
	Parent   *Node
}

// newNode creates a new trie with nothing in it.
func newNode(parent *Node) *Node {
	result := new(Node)
	result.Children = make(map[rune]*Node)
	result.Parent = parent
	return result
}

// Add adds Nodes to the current Node,
func (t *Node) Add(s string) {
	var currentNode *Node = t
	var newS string
	for _, c := range []rune(s) {
		newS += string(c)
		next, ok := currentNode.Children[c]
		if ok {
			currentNode = next
		} else {
			currentNode.Children[c] = newNode(currentNode)
			currentNode = currentNode.Children[c]
		}
	}
	currentNode.Terminal = true
}

// Lookup queries the Node to determine if it contains s.
func (t *Node) Lookup(s string) bool {
	var currentNode *Node = t
	for _, c := range []rune(s) {
		next, ok := currentNode.Children[c]
		if ok {
			currentNode = next
		} else {
			return false
		}
	}
	return currentNode.Terminal
}

// FromFile opens filename and places every line into the trie. It then
// returns the trie.
// If the filename cannot be opened, it returns that error.
func FromFile(filename string) (*Node, error) {
	result := newNode(nil)
	file, err := os.Open(filename)

	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result.Add(strings.TrimSpace(scanner.Text()))
	}
	return result, nil
}
