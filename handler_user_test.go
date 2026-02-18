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

func setupTestState() (*state, command) {
	return &state{
			cfg: &config.Config{
				CurrentUserName: "test_user",
			},
		}, command{
			Name: "Test",
			Args: []string{},
		}
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

func setupTestDB(t *testing.T) (*database.Queries, func()) {
	dbURL := "postgres://gator_user:gator_pass@localhost:5433/gator?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatalf("erro ao conectar no banco de testes: %v", err)
	}

	dbQueries := database.New(db)
	cleanup := func() {
		db.Exec("DELETE FROM users")
		db.Close()
	}
	return dbQueries, cleanup
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
	s, cmd := setupTestState()
	output, err := capturaOutput(func() error {
		return handlerLogin(s, cmd)
	})

	if err == nil {
		t.Logf("output: %s", output)
		t.Fatalf("esperava um erro avisando que nao pode string vazia")
	}
}

func TestHandlerLogin_ValidUserName(t *testing.T) {
	setupTestEnv(t)
	s, cmd := setupTestState()
	cmd.Args = []string{"HenriqueVigato"}

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
	dbQueries, cleanup := setupTestDB(t)
	defer cleanup()

	s := &state{db: dbQueries}
	cmd := command{Name: "register", Args: []string{}}

	err := handlerRegister(s, cmd)
	if err == nil {
		t.Fatalf("esperava erro quando nao ha argumentos")
	}
}

func TestHandlerRegister_UsuarioJaExiste(t *testing.T) {
	dbQueries, cleanup := setupTestDB(t)
	defer cleanup()

	s := &state{db: dbQueries}
	ctx := context.Background()
	_, err := dbQueries.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "Test User",
	})
	if err != nil {
		t.Fatalf("erro ao criar o usuario de test %v", err)
	}

	cmd := command{Name: "register", Args: []string{"Test User"}}
	err = handlerRegister(s, cmd)

	if err == nil {
		t.Fatalf("esperava um erro pois o usuario ja existe")
	}

	if !strings.Contains(err.Error(), cmd.Args[0]) {
		t.Errorf("Esperava que tivesse uma mensagem de sucesso com o nome do usuario criado, mas recebeu %v", err.Error())
	}
}

func TestHandlerRegister_Success(t *testing.T) {
	dbQueries, cleanup := setupTestDB(t)
	defer cleanup()

	s := &state{db: dbQueries}
	cmd := command{Name: "register", Args: []string{"test_user_success"}}

	err := handlerRegister(s, cmd)
	if err != nil {
		t.Fatalf("nao era esperado nenhum erro")
	}

	_, err = s.db.GetUser(context.Background(), "test_user_success")
	if err != nil {
		t.Fatalf("nao era esperado nao encontar nenhum erro")
	}
}
