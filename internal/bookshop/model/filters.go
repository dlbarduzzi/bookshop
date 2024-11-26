package model

import "github.com/dlbarduzzi/bookshop/internal/validator"

type Filters struct {
	Page     int
	PageSize int
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
}

func (f Filters) validatePage(v *validator.Validator) {
	if f.Page < 1 {
		v.AddError("page", "must be greater than 0")
		return
	}
	if f.Page > 10_000_000 {
		v.AddError("page", "cannot be greater than 10 million")
		return
	}
}

func (f *Filters) validatePageSize(v *validator.Validator) {
	if f.PageSize < 1 {
		v.AddError("page_size", "must be greater than 0")
		return
	}
	if f.PageSize > 100 {
		v.AddError("page_size", "cannot be greater than 100")
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
