package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
}

func Read() (*Config, error) {
	var config Config
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter o caminho da home: %w", err)
	}
	gatorConfigJSON := filepath.Join(home, configFileName)
	data, err := os.ReadFile(gatorConfigJSON)
	if err != nil {
		return nil, fmt.Errorf("erro a tentar abrir o arquivo gatorconfig: %w", err)
	}

	if erro := json.Unmarshal(data, &config); erro != nil {
		return nil, fmt.Errorf("erro ao converter o arquivo JSON para struct: %w", erro)
	}

	return &config, nil
}
