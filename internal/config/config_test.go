package config

import (
	"testing"
)

func TestStructConvertion(t *testing.T) {
	config, err := Read()
	if err != nil {
		t.Errorf("%v: ", err)
	}

	if config.DBURL != "postgres://example" {
		t.Logf("Output recebido: %s", config.DBURL)
		t.Errorf("esperava encontrar a string 'postgress://example'")
	}
}
