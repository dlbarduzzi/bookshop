package security

import (
	"regexp"
	"testing"
)

func TestRandomString(t *testing.T) {
	reg := regexp.MustCompile(`[a-zA-Z0-9]+`)
	length := 10
	generated := make([]string, 0, 1000)

	for i := 0; i < 1000; i++ {
		res, err := RandomString(length)
		if err != nil {
			t.Fatal(err)
		}

		if len(res) != length {
			t.Fatalf("expected string length (%d) to be %d; got %d", i, length, len(res))
		}

		if match := reg.MatchString(res); !match {
			t.Fatalf("expected string (%d) to have only [a-zA-Z0-9]+ characters; got %q", i, res)
		}

		for _, str := range generated {
			if str == res {
				t.Fatalf("expected random string (%d) %q not to be repeated; %v", i, res, generated)
			}
		}

		generated = append(generated, res)
	}
}
