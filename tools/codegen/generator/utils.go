package generator

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ompgo-dev/ompgo/tools/codegen/model"
)

// loadAPIGroups reads api.json into API groups.
func loadAPIGroups(path string) (map[string][]model.APIFunction, error) {
	fmt.Printf("Reading %s...\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var groups map[string][]model.APIFunction
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// loadEventGroups reads events.json into event groups.
func loadEventGroups(path string) (map[string][]model.Event, error) {
	fmt.Printf("Reading %s...\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var groups map[string][]model.Event
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}
