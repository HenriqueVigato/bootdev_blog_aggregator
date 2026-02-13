package config

import (
	"os"
	"strings"
	"testing"
)

func setupTestEnv(t *testing.T) string {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	testPath, _ := getConfigFilePath()

	err := os.WriteFile(testPath, []byte(`{"db_url":"postgres://example"}`), 0o644)
	if err != nil {
		t.Fatalf("erro ao escrever no arquivo de teste %v", err)
	}

	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	return tmpDir
}

func TestGetConfigFilePath(t *testing.T) {
	setupTestEnv(t)
	path, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("erro ao obter o path do arquivo config")
	}

	if !strings.Contains(path, ".gatorconfig.json") {
		t.Logf("Output: %v ", path)
		t.Fatalf("esperava encontrar o output: .gatorconfig.json")
	}
}

func TestRead(t *testing.T) {
	setupTestEnv(t)

	config, err := Read()
	if err != nil {
		t.Fatalf("%v: ", err)
	}

	if !strings.Contains(config.DBURL, "postgres://example") {
		t.Logf("Output recebido: %s", config.DBURL)
		t.Fatalf("esperava encontrar a string 'postgress://example'")
	}
}

func TestSetUser(t *testing.T) {
	err := SetUser("Henrique")
	if err != nil {
		t.Fatalf("Erro ao definir o usuario novo %v", err)
	}
	config, err := Read()
	if err != nil {
		t.Fatalf("Erro ao ler o arquivo de config: %v", err)
	}

	if config.Current_user_name != "Henrique" {
		t.Logf("Output recebido: %s", config.Current_user_name)
		t.Fatalf("current user name diferente do esperado")
	}
}
