package pokeapi

import (
	"testing"
	"time"

	cache "github.com/zer0warm/pokedex-repl/internal/pokecache"
)

func TestGetLocationAreas(t *testing.T) {
	cfg := &Config{Cache: cache.NewCache(5 * time.Second)}
	areas, err := cfg.GetLocationAreas(true)
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("areas: %v\n", areas)
	t.Logf("config: %v\n", cfg)
}

func TestGetPokemonInfo(t *testing.T) {
	cases := []string{
		"clefairy",
		"tentacool",
		"tentacruel",
		"pikachu",
	}

	cfg := &Config{Cache: cache.NewCache(5 * time.Second)}

	for _, c := range cases {
		cfg.Args = []string{c}
		info, err := cfg.GetPokemonInfo()
		if err != nil {
			t.Errorf("%v", err)
		}
		t.Log(info)
	}
}
