package model

import (
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func TestFiltersValidate(t *testing.T) {
	t.Parallel()
	v := validator.NewValidator()

	filters := &Filters{
		Page:         1,
		PageSize:     5,
		Sort:         "name",
		SortSafeList: []string{"id", "name"},
	}

	filters.Validate(v)

	if !v.IsValid() {
		t.Error("expected filters validation to succeed")
	}
}

func TestFiltersValidatePage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		page    int
		isValid bool
		wantErr string
	}{
		{
			name:    "negative",
			page:    -1,
			isValid: false,
			wantErr: "Page must be greater than zero.",
		},
		{
			name:    "zero",
			page:    0,
			isValid: false,
			wantErr: "Page must be greater than zero.",
		},
		{
			name:    "to_high",
			page:    20_000_000,
			isValid: false,
			wantErr: "Page cannot be greater than 10 million.",
		},
		{
			name:    "valid",
			page:    10,
			isValid: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &Filters{
				Page: tt.page,
			}

			v := validator.NewValidator()
			f.validatePage(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected filters page validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["page"]
			if !exists && !tt.isValid {
				t.Error("expected page field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected page error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}

func TestFiltersValidatePageSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		pageSize int
		isValid  bool
		wantErr  string
	}{
		{
			name:     "negative",
			pageSize: -1,
			isValid:  false,
			wantErr:  "Page size must be greater than zero.",
		},
		{
			name:     "zero",
			pageSize: 0,
			isValid:  false,
			wantErr:  "Page size must be greater than zero.",
		},
		{
			name:     "to_high",
			pageSize: 200,
			isValid:  false,
			wantErr:  "Page size cannot be greater than 100.",
		},
		{
			name:     "valid",
			pageSize: 10,
			isValid:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &Filters{
				PageSize: tt.pageSize,
			}

			v := validator.NewValidator()
			f.validatePageSize(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected filters page size validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["page_size"]
			if !exists && !tt.isValid {
				t.Error("expected page_size field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected page_size error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}

func TestFiltersValidateSortSafeList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sort     string
		safeList []string
		isValid  bool
		wantErr  string
	}{
		{
			name:     "invalid",
			sort:     "id",
			safeList: []string{"name"},
			isValid:  false,
			wantErr:  "Invalid sort value.",
		},
		{
			name:     "valid",
			sort:     "name",
			safeList: []string{"name"},
			isValid:  true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			f := &Filters{
				Sort:         tt.sort,
				SortSafeList: tt.safeList,
			}

			v := validator.NewValidator()
			f.validateSortSafeList(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected filters sort validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["sort"]
			if !exists && !tt.isValid {
				t.Error("expected sort field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected sort error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}
