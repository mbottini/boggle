package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/mbottini/boggle/graph"
	"github.com/mbottini/boggle/trie"
)

type frame struct {
	Node    *graph.Node
	Trie    *trie.Node
	Visited map[*graph.Node]bool
	Path    map[*graph.Node]bool
	Word    string
}

func createFrame(n *graph.Node,
	t *trie.Node,
	p map[*graph.Node]bool,
	s string) frame {
	var result frame
	result.Node = n
	result.Trie = t
	result.Visited = make(map[*graph.Node]bool)
	if p != nil {
		result.Path = copyMap(p)
	} else {
		result.Path = make(map[*graph.Node]bool)
	}
	result.Word = s
	return result
}

func copyMap(m map[*graph.Node]bool) map[*graph.Node]bool {
	var result map[*graph.Node]bool = make(map[*graph.Node]bool)
	for k := range m {
		result[k] = true
	}
	return result
}

func getWords(n *graph.Node, t *trie.Node) []string {
	result := make(map[string]bool)
	var frameStack []frame
	frameStack = append(frameStack, createFrame(n, t, nil, ""))
	for len(frameStack) > 0 {
		addedFrame := false
		currentFrame := frameStack[len(frameStack)-1]
		for _, candidate := range currentFrame.Node.Connections {
			if !currentFrame.Visited[candidate] && !currentFrame.Path[candidate] {
				t2, found := currentFrame.Trie.Children[candidate.Data]
				if found {
					addedFrame = true
					newWord := currentFrame.Word + string(candidate.Data)
					newFrame := createFrame(candidate, t2, currentFrame.Path, newWord)
					newFrame.Path[currentFrame.Node] = true
					frameStack = append(frameStack, newFrame)
					currentFrame.Visited[candidate] = true
					if t2.Terminal && len(newFrame.Word) >= 4 {
						result[newFrame.Word] = true
					}
					break
				}
			}
		}
		// If we didn't find anything from the current frame, we pop it off the
		// stack.
		if !addedFrame {
			frameStack = frameStack[:len(frameStack)-1]
		}
	}

	var keys []string
	for k := range result {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	t, err := trie.FromFile("/home/mike/dictionary.txt")
	if err != nil {
		fmt.Println("Couldn't open file!")
		os.Exit(1)
	}

	f, err := graph.FromFile("/home/mike/boggle.txt")
	if err != nil {
		fmt.Println("Couldn't open Boggle field!")
		os.Exit(1)
	}

	result := getWords(f, t)
	sort.Strings(result)
	for _, w := range result {
		fmt.Println(w)
	}
	fmt.Printf("%d words found\n", len(result))
}
