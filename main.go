package main

type TrieNode struct {
	children map[rune]*TrieNode
	meaning  string
	isWord   bool
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
}
