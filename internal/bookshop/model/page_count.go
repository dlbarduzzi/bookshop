package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidPageCountFormat = errors.New("invalid page count format")

type PageCount int32

func (p PageCount) MarshalJSON() ([]byte, error) {
	pages := fmt.Sprintf("%d pages", p)
	quote := strconv.Quote(pages)
	return []byte(quote), nil
}

func (p *PageCount) UnmarshalJSON(value []byte) error {
	pages, err := strconv.Unquote(string(value))
	if err != nil {
		return ErrInvalidPageCountFormat
	}

	parts := strings.Split(pages, " ")

	if len(parts) != 2 || parts[1] != "pages" {
		return ErrInvalidPageCountFormat
	}

	count, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidPageCountFormat
	}

	*p = PageCount(count)

	return nil
}
