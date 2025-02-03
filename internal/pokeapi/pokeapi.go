package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	cache "github.com/zer0warm/pokedex-repl/internal/pokecache"
)

type Area struct {
	Name string `json:"name"`
}

type PokemonInfo struct {
	Name    string `json:"name"`
	BaseEXP int    `json:"base_experience"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	Stats   []struct {
		Base int `json:"base_stat"`
		Stat struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func (info PokemonInfo) String() string {
	result := ""
	result += fmt.Sprintln("Name:", info.Name)
	result += fmt.Sprintln("Height:", info.Height)
	result += fmt.Sprintln("Weight:", info.Weight)
	result += fmt.Sprintln("Stats:")
	for _, stat := range info.Stats {
		result += fmt.Sprintf("- %s: %d\n", stat.Stat.Name, stat.Base)
	}
	result += fmt.Sprintln("Types:")
	for _, typ := range info.Types {
		result += fmt.Sprintf("- %s\n", typ.Type.Name)
	}
	return result[:len(result)-1]
}

type Config struct {
	Next     string
	Previous string
	Args     []string
	Cache    *cache.Cache
	Pokedex  map[string]PokemonInfo
}

var (
	ErrNoNext = errors.New("no next entries")
	ErrNoPrev = errors.New("no prev entries")
)

// Get location areas with direction (forward or not)
func (cfg *Config) GetLocationAreas(forward bool) ([]Area, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

	// Err when running outbounds
	if forward && cfg.Next == "" && cfg.Previous != "" {
		return nil, ErrNoNext
	}
	if !forward && cfg.Previous == "" {
		return nil, ErrNoPrev
	}

	// Update endpoint URL when possible
	if forward && cfg.Next != "" {
		endpoint = cfg.Next
	}
	if !forward && cfg.Previous != "" {
		endpoint = cfg.Previous
	}

	body, err := cfg.cacheGet(endpoint)
	if err != nil {
		return nil, err
	}

	data := struct {
		Next     string `json:"next,omitempty"`
		Previous string `json:"previous,omitempty"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("while decoding response location-area: %w", err)
	}

	cfg.Next = data.Next
	cfg.Previous = data.Previous

	areas := []Area{}
	for _, r := range data.Results {
		areas = append(areas, Area{Name: r.Name})
	}

	return areas, nil
}

// Get pokemons can be encountered in an area with their name
// Name is supplied in config.Args
func (cfg *Config) GetAreaPokemons() ([]string, error) {
	endpoint := fmt.Sprintf(
		"https://pokeapi.co/api/v2/location-area/%s?offset=0&limit=20",
		cfg.Args[0],
	)

	body, err := cfg.cacheGet(endpoint)
	if err != nil {
		return nil, err
	}

	data := struct {
		PokemonEncounters []struct {
			Pokemon struct {
				Name string `json:"name"`
			} `json:"pokemon"`
		} `json:"pokemon_encounters"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("while decoding JSON data: %w", err)
	}

	pokemons := []string{}
	for _, pokemon := range data.PokemonEncounters {
		pokemons = append(pokemons, pokemon.Pokemon.Name)
	}

	return pokemons, nil
}

// Get pokemon information
// For now, only their name (echo from request) and their base EXP
func (cfg *Config) GetPokemonInfo() (PokemonInfo, error) {
	endpoint := fmt.Sprintf(
		"https://pokeapi.co/api/v2/pokemon/%s",
		cfg.Args[0],
	)

	body, err := cfg.cacheGet(endpoint)
	if err != nil {
		return PokemonInfo{}, fmt.Errorf("while GET: %w", err)
	}

	info := PokemonInfo{}
	if err := json.Unmarshal(body, &info); err != nil {
		return PokemonInfo{}, fmt.Errorf("while decoding JSON: %w", err)
	}

	return info, nil
}

// Use cache, otherwise add to cache
func (cfg *Config) cacheGet(endpoint string) ([]byte, error) {
	var body []byte

	if value, ok := cfg.Cache.Get(endpoint); ok {
		log.Printf("Cache: HIT - %s\n", endpoint)
		body = value
	} else {
		log.Printf("Cache: MISS - %s\n", endpoint)

		res, err := http.Get(endpoint)
		if err != nil {
			return nil, fmt.Errorf("while GET location-area: %w", err)
		}
		defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("while consuming response: %w", err)
		}

		cfg.Cache.Add(endpoint, body)
	}

	return body, nil
}
