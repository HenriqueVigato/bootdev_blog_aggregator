package main

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/config"
	"github.com/HenriqueVigato/bootdev_blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func setupTestState(t *testing.T) (*state, command, string) {
	db := setupTestDB(t)
	return &state{
			db: db,
			cfg: &config.Config{
				DBURL:           "",
				CurrentUserName: "test_user",
			},
		}, command{
			Name: "Test",
			Args: []string{},
		},
		setupTestEnv(t)
}

func setupTestEnv(t *testing.T) string {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	home, _ := os.UserHomeDir()

	testPath := filepath.Join(home, ".gatorconfig.json")

	err := os.WriteFile(testPath, []byte(`{"db_url":"postgres://example"}`), 0o644)
	if err != nil {
		t.Fatalf("erro ao escrever no arquivo de teste %v", err)
	}

	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	return tmpDir
}

func setupTestDB(t *testing.T) *database.Queries {
	dbURL := "postgres://gator_user:gator_pass@localhost:5433/gator?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("erro ao conectar no banco de testes: %v", err)
	}

	dbQueries := database.New(db)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users")
		db.Close()
	})
	return dbQueries
}

func capturaOutput(fn func() error) (string, error) {
	old := os.Stdout
	defer func() {
		os.Stdout = old
	}()

	r, w, _ := os.Pipe()
	os.Stdout = w

	err := fn()

	w.Close()

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestHandlerLogin_emptyUserName(t *testing.T) {
	setupTestEnv(t)
	s, cmd, _ := setupTestState(t)
	output, err := capturaOutput(func() error {
		return handlerLogin(s, cmd)
	})

	if err == nil {
		t.Logf("output: %s", output)
		t.Fatalf("esperava um erro avisando que nao pode string vazia")
	}
}

func TestHandlerLogin_UnregisteredUserName(t *testing.T) {
	setupTestEnv(t)
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"Unregistered user"}

	_, err := capturaOutput(func() error {
		return handlerLogin(s, cmd)
	})

	if err == nil {
		t.Errorf("Deveria retorar um erro pois o user nao esta cadastrado no banco de dados")
	}

	if !strings.Contains(err.Error(), "ja estar criado") {
		t.Logf("unregistered user: %v", err.Error())
		t.Errorf("deveria aparecer um erro dizendo que o usuario deveria estar cadastrado")
	}
}

func TestHandlerLogin_ValidUserName(t *testing.T) {
	setupTestEnv(t)
	s, cmd, _ := setupTestState(t)
	cmd.Args = []string{"Test valid user"}

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Test valid user",
	})
	if err != nil {
		t.Fatalf("erro ao criar o usuario de test %v", err)
	}

	output, err := capturaOutput(func() error {
		return handlerLogin(s, cmd)
	})
	if err != nil {
		t.Fatalf("erro %v:", err)
	}

	cfg, err := config.Read()
	if err != nil {
		t.Logf("s.cfg.CurrentUserName %s", cfg.CurrentUserName)
		t.Fatalf("CurrentUserName diferento do nome definido")
	}

	if !strings.Contains(output, "been set") {
		t.Logf("output: %v", output)
		t.Fatalf("Esperava a mensagem de sucesso")
	}
}

func TestHandlerRegister_FaltaArgumento(t *testing.T) {
	s, cmd, _ := setupTestState(t)

	cmd.Name = "register"
	cmd.Args = []string{}

	err := handlerRegister(s, cmd)
	if err == nil {
		t.Fatalf("esperava erro quando nao ha argumentos")
	}
}

func TestHandlerRegister_UsuarioJaExiste(t *testing.T) {
	s, cmd, _ := setupTestState(t)

	ctx := context.Background()
	_, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Test User",
	})
	if err != nil {
		t.Fatalf("erro ao criar o usuario de test %v", err)
	}

	cmd.Name = "register"
	cmd.Args = []string{"Test User"}

	err = handlerRegister(s, cmd)

	if err == nil {
		t.Fatalf("esperava um erro pois o usuario ja existe")
	}

	if !strings.Contains(err.Error(), cmd.Args[0]) {
		t.Errorf("Esperava que tivesse uma mensagem de sucesso com o nome do usuario criado, mas recebeu %v", err.Error())
	}
}

func TestHandlerRegister_Success(t *testing.T) {
	s, cmd, _ := setupTestState(t)

	cmd.Name = "register"
	cmd.Args = []string{"test_user_success"}

	err := handlerRegister(s, cmd)
	if err != nil {
		t.Fatalf("nao era esperado nenhum erro: %v", err)
	}

	_, err = s.db.GetUser(context.Background(), "test_user_success")
	if err != nil {
		t.Fatalf("nao era esperado nao encontar nenhum erro")
	}
}

func TestHandlerReset(t *testing.T) {
	s, cmd, _ := setupTestState(t)

	cmd.Name = "register"
	cmd.Args = []string{"test_user_delete"}

	err := handlerRegister(s, cmd)
	if err != nil {
		t.Fatalf("nao era esperado nenhum erro aqui %v", err)
	}

	user, err := s.db.GetUser(context.Background(), "test_user_delete")
	if err != nil {
		t.Logf("user: %v", user)
		t.Fatalf("nao era esperado nenhum erro agoria na hora de ver se o usuario teste foi criado")
	}

	err = handlerReset(s, cmd)
	if err != nil {
		t.Logf("resetUser err: %v", err)
		t.Fatalf("nao era esperado nenhum erro aqui")
	}

	_, err = s.db.GetUser(context.Background(), "test_user_delete")
	if err == nil {
		t.Logf("err %v", err)
		t.Fatalf("era esperado um erro pois o usuario foi excluido")
	}
}
