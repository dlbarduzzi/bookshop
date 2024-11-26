package bookshop

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func (b *Bookshop) readString(q url.Values, k string, defaultValue string) string {
	s := q.Get(k)
	if s == "" {
		return defaultValue
	}
	return s
}

func (b *Bookshop) readCSV(q url.Values, k string, defaultValue []string) []string {
	csv := q.Get(k)
	if csv == "" {
		return defaultValue
	}
	return strings.Split(csv, ",")
}

func (b *Bookshop) readInt(q url.Values, k string, defaultValue int, v *validator.Validator) int {
	s := q.Get(k)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(k, "must be an integer value")
		return defaultValue
	}
	return i
}
