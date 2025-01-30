package main

import (
	"fmt"
	"os"

	api "github.com/zer0warm/pokedex-repl/internal/pokeapi"
)

// cliCommand represents a CLI command sent to the Pokedex.
type cliCommand struct {
	name     string
	desc     string
	callback func(cfg *api.Config) error
}

func listCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:     "exit",
			desc:     "Exit the pokedex",
			callback: commandExit,
		},
		"help": {
			name:     "help",
			desc:     "Display a help message",
			callback: commandHelp,
		},
		"map": {
			name:     "map",
			desc:     "List 20 location areas",
			callback: commandMap,
		},
	}
}

func commandExit(cfg *api.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *api.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range listCommands() {
		fmt.Printf("%-20s%s\n", cmd.name, cmd.desc)
	}
	return nil
}

func commandMap(cfg *api.Config) error {
	areas, err := cfg.GetLocationAreas()
	if err != nil {
		return fmt.Errorf("while getting location-areas: %w", err)
	}

	for _, area := range areas {
		fmt.Println(area.Name)
	}

	return nil
}
