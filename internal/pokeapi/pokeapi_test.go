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
