package model

import (
	"fmt"
	"strconv"
)

type PageCount int32

func (p PageCount) MarshalJSON() ([]byte, error) {
	pages := fmt.Sprintf("%d pages", p)
	quote := strconv.Quote(pages)
	return []byte(quote), nil
}
