package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	var config Config
	gatorConfigJSON, erro := getConfigFilePath()
	if erro != nil {
		return nil, fmt.Errorf("erro ao obter o caminho do arquivo config %w", erro)
	}

	data, err := os.ReadFile(gatorConfigJSON)
	if err != nil {
		return nil, fmt.Errorf("erro a tentar abrir o arquivo gatorconfig: %w", err)
	}

	if erro := json.Unmarshal(data, &config); erro != nil {
		return nil, fmt.Errorf("erro ao converter o arquivo JSON para struct: %w", erro)
	}

	return &config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	const configFileName = ".gatorconfig.json"
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter o caminho da home: %w", err)
	}

	return filepath.Join(home, configFileName), nil
}

func write(config *Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	dadosByte, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, dadosByte, 0o664)
}
