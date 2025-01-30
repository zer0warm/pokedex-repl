package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Area struct {
	Name string `json:"name"`
}

type Config struct {
	Next     string
	Previous string
}

var (
	ErrNoNext = errors.New("no next entries")
	ErrNoPrev = errors.New("no prev entries")
)

// Get location areas with direction (forward or not)
func (cfg *Config) GetLocationAreas(forward bool) ([]Area, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"

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

	res, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("while GET location-area: %w", err)
	}
	defer res.Body.Close()

	data := struct {
		Next     string `json:"next,omitempty"`
		Previous string `json:"previous,omitempty"`
		Results  []struct {
			Name string `json:"name"`
		} `json:"results"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
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
