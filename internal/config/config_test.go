package config

import (
	"testing"
)

func TestGetConfigFilePath(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Errorf("erro ao obter o path do arquivo config")
	}

	if path != "/home/henrique/.gatorconfig.json" {
		t.Logf("Output: %v ", path)
		t.Errorf("esperava encontrar o output: /home/henrique/.gatorconfig.json")
	}
}

func TestRead(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Errorf("%v: ", err)
	}

	if config.DBURL != "postgres://example" {
		t.Logf("Output recebido: %s", config.DBURL)
		t.Errorf("esperava encontrar a string 'postgress://example'")
	}
}

func TestSetUser(t *testing.T) {
	err := SetUser("Henrique")
	if err != nil {
		t.Errorf("Erro ao definir o usuario novo %v", err)
	}
	config, err := Read()
	if err != nil {
		t.Errorf("Erro ao ler o arquivo de config: %v", err)
	}

	if config.Current_user_name != "Henrique" {
		t.Logf("Output recebido: %s", config.Current_user_name)
		t.Errorf("current user name diferente do esperado")
	}
}
