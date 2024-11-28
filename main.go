package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// File path for the dictionary
const dictionaryFile = "word_dictionary.txt"

// Function to add a word to the dictionary
func addWord(word, meaning string) {
	file, err := os.OpenFile(dictionaryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening dictionary file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s:%s\n", word, meaning))
	if err != nil {
		fmt.Println("Error writing to dictionary file:", err)
		return
	}
	fmt.Printf("Word '%s' added to the dictionary.\n", word)
}

// Function to load the dictionary into memory
func loadDictionary() map[string]string {
	dictionary := make(map[string]string)

	file, err := os.Open(dictionaryFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Dictionary file not found. A new one will be created.")
		} else {
			fmt.Println("Error opening dictionary file:", err)
		}
		return dictionary
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			dictionary[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading dictionary file:", err)
	}

	return dictionary
}

// Function to query a word in the dictionary
func queryWord(word string) string {
	dictionary := loadDictionary()
	if meaning, exists := dictionary[word]; exists {
		return meaning
	}
	return "Word not found in the dictionary."
}

func main() {
	fmt.Println("Welcome to the Word Dictionary!")

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
			addWord(word, meaning)
		case 2:
			var word string
			fmt.Print("Enter the word to query: ")
			fmt.Scanln(&word)
			fmt.Println("Meaning:", queryWord(word))
		case 3:
			fmt.Println("Exiting. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
