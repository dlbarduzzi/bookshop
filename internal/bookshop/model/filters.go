package model

import (
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f *Filters) Validate(v *validator.Validator) {
	f.validatePage(v)
	f.validatePageSize(v)
	f.validateSortSafeList(v)
}

func (f *Filters) validatePage(v *validator.Validator) {
	page := f.Page

	if page < 1 {
		v.AddError("page", "Page must be greater than zero.")
		return
	}

	if page > 10_000_000 {
		v.AddError("page", "Page cannot be greater than 10 million.")
		return
	}
}

func (f *Filters) validatePageSize(v *validator.Validator) {
	pageSize := f.PageSize

	if pageSize < 1 {
		v.AddError("page_size", "Page size must be greater than zero.")
		return
	}

	if pageSize > 100 {
		v.AddError("page_size", "Page size cannot be greater than 100.")
		return
	}
}

func (f *Filters) validateSortSafeList(v *validator.Validator) {
	if !validator.ValueInList(f.Sort, f.SortSafeList...) {
		v.AddError("sort", "Invalid sort value.")
		return
	}
}
