package pokeapi

import "testing"

func TestGetLocationAreas(t *testing.T) {
	cfg := &Config{}
	areas, err := cfg.GetLocationAreas()
	if err != nil {
		t.Errorf("%v", err)
	}
	t.Logf("areas: %v\n", areas)
	t.Logf("config: %v\n", cfg)
}
