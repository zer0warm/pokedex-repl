package main

import (
	"errors"
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
		"mapb": {
			name:     "mapb",
			desc:     "List previous 20 location areas",
			callback: commandMapb,
		},
		"explore": {
			name:     "explore",
			desc:     "List pokemons can be encountered in an area",
			callback: commandExplore,
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
	areas, err := cfg.GetLocationAreas(true)
	if err != nil {
		if errors.Is(err, api.ErrNoNext) {
			fmt.Println("You are on the last page!")
			return nil
		}
		return fmt.Errorf("while getting location-areas: %w", err)
	}

	for _, area := range areas {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(cfg *api.Config) error {
	areas, err := cfg.GetLocationAreas(false)
	if err != nil {
		if errors.Is(err, api.ErrNoPrev) {
			fmt.Println("You are on the first page!")
			return nil
		}
		return fmt.Errorf("while getting location-areas: %w", err)
	}

	for _, area := range areas {
		fmt.Println(area.Name)
	}

	return nil
}

func commandExplore(cfg *api.Config) error {
	if len(cfg.Args) != 1 {
		return fmt.Errorf("must supply 1 location")
	}

	pokemons, err := cfg.GetAreaPokemons()
	if err != nil {
		return fmt.Errorf(
			"while getting list of pokemons encountered in %s: %w",
			cfg.Args[0], err)
	}

	for _, pokemon := range pokemons {
		fmt.Println(pokemon)
	}

	return nil
}
