package helpers

import (
	"encoding/json"
	"os"
)

type FileCred struct {
	Service  string
	Username string
	Password string
}

func readJsonFile(f string) (map[string][]map[string]string, error) {
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

func writeJsonFile(f string, data map[string][]map[string]string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(f, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

func ReadCredsFile(f string) ([]FileCred, error) {
	jsonData, err := readJsonFile(f)
	if err != nil {
		return nil, err
	}
	var creds []FileCred
	for _, credMap := range jsonData["creds_list"] {
		cred := FileCred{
			Service:  credMap["service"],
			Username: credMap["username"],
			Password: credMap["password"],
		}
		creds = append(creds, cred)
	}
	return creds, nil
}

func WriteCredsFile(f string, creds []FileCred) error {
	var jsonData = map[string][]map[string]string{"creds_list": {}}
	for _, cred := range creds {
		credMap := map[string]string{
			"service":  cred.Service,
			"username": cred.Username,
			"password": cred.Password,
		}
		jsonData["creds_list"] = append(jsonData["creds_list"], credMap)
	}
	return writeJsonFile(f, jsonData)
}
