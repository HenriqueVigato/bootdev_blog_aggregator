package main

import (
	"strings"
	"testing"
)

func TestAgg(t *testing.T) {
	s, cmd, _ := setupTestState(t)
	output, err := capturaOutput(func() error {
		return agg(s, cmd)
	})
	if err != nil {
		t.Fatalf("erro ao executar o comando agg, %v", err)
	}

	if !strings.Contains(output, "Lane's Blog") {
		t.Logf("output \n %v", output)
		t.Errorf("output diferente do esperado")
	}
}
