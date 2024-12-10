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

		inputLine := strings.Trim(strings.ToLower(scanner.Text()), " \t")
		command := strings.Fields(inputLine)[0]
		fmt.Println("Your command was:", command)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
