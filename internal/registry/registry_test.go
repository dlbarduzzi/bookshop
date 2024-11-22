package registry

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNewRegistry(t *testing.T) {
	SetConfigPath("/path/not/found")
	SetConfigName("test")
	SetConfigType("yaml")

	_, err := NewRegistry()
	if err == nil {
		t.Fatal("expected error not to be nil")
	}

	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		t.Fatalf(
			"expected error to be %v; got %v",
			err.(viper.ConfigFileNotFoundError).Error(),
			err,
		)
	}

	SetConfigPath(".")
	SetConfigName("test")
	SetConfigType("yaml")

	reg, err := NewRegistry()
	if err != nil {
		t.Fatalf("expected error to be nil; got %v", err)
	}

	env := reg.GetString("TEST_VAR")
	if env != "test-value" {
		t.Fatal(err)
	}
}
