package pokeapi

import (
	"encoding/json"
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

func (cfg *Config) GetLocationAreas() ([]Area, error) {
	endpoint := "https://pokeapi.co/api/v2/location-area/"
	if cfg.Next != "" {
		endpoint = cfg.Next
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
