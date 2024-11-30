package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// TrieNode represents a node in the Trie
type TrieNode struct {
	children map[rune]*TrieNode
	meaning  string
	isWord   bool
}

// Trie represents the Trie data structure
type Trie struct {
	root *TrieNode
	file *os.File // File handle for saving
}

// NewTrie creates a new Trie with a file for persistence
func NewTrie(filePath string) (*Trie, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &Trie{
		root: &TrieNode{children: make(map[rune]*TrieNode)},
		file: file,
	}, nil
}

// AddWord adds a word and its meaning to the Trie and saves it to the file
func (t *Trie) AddWord(word, meaning string) error {
	node := t.root
	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[ch]
	}
	node.isWord = true
	node.meaning = meaning

	// Save to file immediately
	_, err := t.file.WriteString(fmt.Sprintf("%s:%s\n", word, meaning))
	if err != nil {
		return fmt.Errorf("failed to save word: %w", err)
	}
	return nil
}

// QueryWord searches for a word in the Trie and returns its meaning
func (t *Trie) QueryWord(word string) string {
	node := t.root
	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			return "Word not found in the dictionary."
		}
		node = node.children[ch]
	}
	if node.isWord {
		return node.meaning
	}
	return "Word not found in the dictionary."
}

// LoadTrie loads words and meanings from the file into the Trie
func LoadTrie(filePath string) (*Trie, error) {
	trie, err := NewTrie(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return trie, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			trie.AddWord(parts[0], parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return trie, nil
}

// Close closes the file associated with the Trie
func (t *Trie) Close() error {
	return t.file.Close()
}

// Main program
func main() {
	const dictionaryFile = "word_dictionary_trie.txt"

	// Load Trie from file
	trie, err := LoadTrie(dictionaryFile)
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}
	defer trie.Close()

	fmt.Println("Welcome to the Trie Word Dictionary!")
	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Add Word")
		fmt.Println("2. Query Word")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		_, err := fmt.Scanf("%d\n", &choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			var word, meaning string
			fmt.Print("Enter the word: ")
			fmt.Scanln(&word)
			fmt.Print("Enter the meaning: ")
			meaning, _ = bufio.NewReader(os.Stdin).ReadString('\n')
			meaning = strings.TrimSpace(meaning)
			if err := trie.AddWord(word, meaning); err != nil {
				fmt.Println("Error adding word:", err)
			} else {
				fmt.Printf("Word '%s' added to the dictionary.\n", word)
			}
		case 2:
			var word string
			fmt.Print("Enter the word to query: ")
			fmt.Scanln(&word)
			fmt.Println("Meaning:", trie.QueryWord(word))
		case 3:
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
