package model

import (
	"strings"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
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

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}
