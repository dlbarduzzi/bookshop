package bookshop

import "testing"

func TestPort(t *testing.T) {
	t.Parallel()

	bs := &Bookshop{
		config: &Config{
			Port: 9191,
		},
	}

	wantPort := bs.config.Port

	if bs.Port() != wantPort {
		t.Errorf("expected bookshop port to be %d; got %d", wantPort, bs.Port())
	}
}
