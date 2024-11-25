package bookshop

import (
	"net/url"
	"strconv"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

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
