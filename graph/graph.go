package graph

import (
	"bufio"
	"os"
	"strings"
)

// Node is a vertex on the graph. It holds a Unicode rune and a list of
// connections to other Nodes.
type Node struct {
	Data        rune
	Root        bool
	Connections []*Node
}

// Connect creates an edge between two vertices.
func (n *Node) Connect(n2 *Node) {
	n.Connections = append(n.Connections, n2)
	n2.Connections = append(n2.Connections, n)
}

type nodePair struct {
	n1 *Node
	n2 *Node
}

type rowPair struct {
	r1 []*Node
	r2 []*Node
}

func zipNodes(n1s []*Node, n2s []*Node) []nodePair {
	var minLen int
	var result []nodePair
	if len(n1s) < len(n2s) {
		minLen = len(n1s)
	} else {
		minLen = len(n2s)
	}
	for i := 0; i < minLen; i++ {
		result = append(result, nodePair{n1: n1s[i], n2: n2s[i]})
	}
	return result
}

func pairNodes(ns []*Node) []nodePair {
	return zipNodes(ns, ns[1:])
}

func zipRows(n1s [][]*Node, n2s [][]*Node) []rowPair {
	var minLen int
	var result []rowPair
	if len(n1s) < len(n2s) {
		minLen = len(n1s)
	} else {
		minLen = len(n2s)
	}
	result = make([]rowPair, minLen)
	for i := 0; i < minLen; i++ {
		result = append(result, rowPair{r1: n1s[i], r2: n2s[i]})
	}
	return result
}

func pairRows(ns [][]*Node) []rowPair {
	return zipRows(ns, ns[1:])
}

func newNode(r rune) *Node {
	var result = new(Node)
	result.Data = r
	return result
}

func newRootNode() *Node {
	var result *Node = new(Node)
	result.Root = true
	return result
}

func createRow(s string) []*Node {
	var result []*Node
	for _, r := range []rune(s) {
		result = append(result, newNode(r))
	}
	for _, pair := range pairNodes(result) {
		pair.n1.Connect(pair.n2)
	}
	return result
}

func connectRows(n1s []*Node, n2s []*Node) {
	for i := 0; i < len(n1s); i++ {
		// We always connect with nodes directly adjacent to us.
		n1s[i].Connect(n2s[i])
		// Diagonal down to the left.
		if i > 0 {
			n1s[i].Connect(n2s[i-1])
		}
		// Diagonal down to the right.
		if i < len(n1s)-1 {
			n1s[i].Connect(n2s[i+1])
		}
	}
}

func createField(ss []string) [][]*Node {
	var result [][]*Node
	for _, s := range ss {
		result = append(result, createRow(s))
	}
	for _, p := range pairRows(result) {
		connectRows(p.r1, p.r2)
	}
	return result
}

// FromFile takes a filename, opens it, parses the lines of runes into a Field,
// and returns it.
func FromFile(filename string) (*Node, error) {
	var result *Node = newRootNode()
	var ss []string
	ss = make([]string, 0)
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ss = append(ss, strings.TrimSpace(scanner.Text()))
	}
	allNodes := createField(ss)
	for _, row := range allNodes {
		for _, node := range row {
			result.Connect(node)
		}
	}
	return result, nil
}
