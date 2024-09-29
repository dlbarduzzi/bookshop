package bookshop

import "testing"

func TestParseConfig(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Port: 0,
	}

	_, err := cfg.parseConfig()
	want := "env variable BS_APP_PORT is missing or invalid"

	if err == nil || err.Error() != want {
		t.Errorf("expected error to be %v; got %v", want, err)
	}

	cfg.Port = 8080

	_, err = cfg.parseConfig()
	if err != nil {
		t.Errorf("expected error to be nil; got %v", err)
	}
}
