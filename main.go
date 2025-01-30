package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex> ")

		if !scanner.Scan() {
			break
		}

		// Do nothing when issuing line breaks
		if len(scanner.Text()) > 0 {
			inputWords := cleanInput(scanner.Text())
			fmt.Println("Your command was:", inputWords[0])
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Split user input into words (whitespace-separated)
func cleanInput(text string) []string {
	return strings.Fields(strings.Trim(strings.ToLower(text), " \t"))
}
