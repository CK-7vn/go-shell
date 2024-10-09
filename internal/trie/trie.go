package trie

import (
	"os"
	"strings"
)

type Node struct {
	Char     string
	Children map[rune]*Node
	isEnd    bool
}

type Trie struct {
	RootNode *Node
}

func NewTrie() *Trie {
	return &Trie{RootNode: &Node{Children: make(map[rune]*Node)}}
}

func (t *Trie) Insert(word string) {
	node := t.RootNode
	for _, ch := range word {
		if _, exists := node.Children[ch]; !exists {
			node.Children[ch] = &Node{Children: make(map[rune]*Node)}
		}
		node = node.Children[ch]
	}
	node.isEnd = true
}

func (t *Trie) Search(prefix string) []string {
	node := t.RootNode
	for _, ch := range prefix {
		if _, exists := node.Children[ch]; !exists {
			return []string{}
		}
		node = node.Children[ch]
	}
	return t.CollectWords(node, prefix)
}

func (t *Trie) CollectWords(node *Node, prefix string) []string {
	result := []string{}
	if node.isEnd {
		result = append(result, prefix)
	}
	for ch, child := range node.Children {
		result = append(result, t.CollectWords(child, prefix+string(ch))...)
	}
	return result
}

func PopulateTrieFromPath(trie *Trie, path string) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}
		for _, file := range files {
			// filePath := filepath.Join(path, file.Name())
			trie.Insert(file.Name())
		}
	}
}

func isExecutable(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}
