package helpers

import (
	"encoding/json"
	"os"
)

func ReadJsonFile(f string) (map[string][]map[string]string, error) {
	data, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}

	var result map[string][]map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func WriteJsonFile(f string, data map[string][]map[string]string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(f, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}
