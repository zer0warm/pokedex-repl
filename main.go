package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	api "github.com/zer0warm/pokedex-repl/internal/pokeapi"
	cache "github.com/zer0warm/pokedex-repl/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	apiConfig := api.Config{
		Cache: cache.NewCache(60 * time.Second),
	}
	for {
		fmt.Print("Pokedex> ")

		if !scanner.Scan() {
			break
		}

		// Do nothing when issuing line breaks
		if len(scanner.Text()) > 0 {
			inputWords := cleanInput(scanner.Text())
			command := inputWords[0]
			apiConfig.Args = inputWords[1:]
			if cmd, ok := listCommands()[command]; ok {
				if err := cmd.callback(&apiConfig); err != nil {
					fmt.Errorf("while running callback: %w\n", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
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
