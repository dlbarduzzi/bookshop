package model

import "testing"

func TestPageCountMarshalJSON(t *testing.T) {
	t.Parallel()

	var pageCount PageCount = 10

	res, err := pageCount.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	want := `"10 pages"`

	if string(res) != want {
		t.Errorf("expected page count to be %s; got %s", want, string(res))
	}
}

func TestPageCountUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pages     string
		pageCount PageCount
		err       error
	}{
		{
			name:      "not_quoted",
			pages:     "",
			pageCount: 0,
			err:       ErrInvalidPageCountFormat,
		},
		{
			name:      "no_pages_string",
			pages:     `"100"`,
			pageCount: 0,
			err:       ErrInvalidPageCountFormat,
		},
		{
			name:      "no_count_value",
			pages:     `"pages"`,
			pageCount: 0,
			err:       ErrInvalidPageCountFormat,
		},
		{
			name:      "bad_format",
			pages:     `"100 pages here"`,
			pageCount: 0,
			err:       ErrInvalidPageCountFormat,
		},
		{
			name:      "no_integer",
			pages:     `"abc pages"`,
			pageCount: 0,
			err:       ErrInvalidPageCountFormat,
		},
		{
			name:      "valid",
			pages:     `"100 pages"`,
			pageCount: 100,
			err:       nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var pageCount PageCount

			err := pageCount.UnmarshalJSON([]byte(tt.pages))
			if err != tt.err {
				t.Errorf("expected page count error to be %v; got %v", tt.err, err)
			}

			if pageCount != tt.pageCount {
				t.Errorf("expected page count to be %d; got %d", tt.pageCount, pageCount)
			}
		})
	}
}
