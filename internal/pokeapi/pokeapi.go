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

type Config struct {
	Next     string
	Previous string
	Cache    *cache.Cache
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
