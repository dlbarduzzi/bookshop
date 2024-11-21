package bookshop

import "testing"

func TestConfigParse(t *testing.T) {
	t.Parallel()

	config := &Config{}

	_, err := config.parse()
	wantErr := "invalid bookshop '0' port number"

	if err == nil || err.Error() != wantErr {
		t.Errorf("expected error to be %v; got %v", wantErr, err)
	}

	config.Port = 8000

	_, err = config.parse()
	if err != nil {
		t.Fatalf("expected error to be nil; got %v", err)
	}
}
