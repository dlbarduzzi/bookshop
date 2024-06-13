package bookshop

import "testing"

func TestRoutes(t *testing.T) {
	t.Parallel()
	bs := newTestBookshop(t)
	if bs.Routes() == nil {
		t.Errorf("expected bookshop routes not to be nil")
	}
}
