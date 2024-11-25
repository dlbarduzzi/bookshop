package model

import (
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func TestValidateFilters(t *testing.T) {
	t.Parallel()

	f := &Filters{
		Page:     1,
		PageSize: 10,
	}

	v := validator.NewValidator()
	f.Validate(v)

	if !v.IsValid() {
		t.Error("expected filters validation to be successful")
	}
}

func TestValidatePage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		page    int
		isValid bool
		wantErr string
	}{
		{
			name:    "less_than_zero",
			page:    -1,
			isValid: false,
			wantErr: "must be greater than 0",
		},
		{
			name:    "zero",
			page:    0,
			isValid: false,
			wantErr: "must be greater than 0",
		},
		{
			name:    "more_than_limit",
			page:    20_000_000,
			isValid: false,
			wantErr: "cannot be greater than 10 million",
		},
		{
			name:    "success",
			page:    1,
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := &Filters{Page: tc.page}
			v := validator.NewValidator()
			f.validatePage(v)

			if v.IsValid() != tc.isValid {
				t.Errorf(
					"expected page filter validation to be %t; got %t",
					tc.isValid,
					v.IsValid(),
				)
			}

			err, ok := v.Errors["page"]
			if !ok && !tc.isValid {
				t.Fatal("expected page field to exist in validation error")
			}

			if err != tc.wantErr {
				t.Errorf("expected page error to be %s; got %s", tc.wantErr, err)
			}
		})
	}
}

func TestValidatePageSize(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		pageSize int
		isValid  bool
		wantErr  string
	}{
		{
			name:     "less_than_zero",
			pageSize: -1,
			isValid:  false,
			wantErr:  "must be greater than 0",
		},
		{
			name:     "zero",
			pageSize: 0,
			isValid:  false,
			wantErr:  "must be greater than 0",
		},
		{
			name:     "more_than_limit",
			pageSize: 200,
			isValid:  false,
			wantErr:  "cannot be greater than 100",
		},
		{
			name:     "success",
			pageSize: 100,
			isValid:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f := &Filters{PageSize: tc.pageSize}
			v := validator.NewValidator()
			f.validatePageSize(v)

			if v.IsValid() != tc.isValid {
				t.Errorf(
					"expected page_size filter validation to be %t; got %t",
					tc.isValid,
					v.IsValid(),
				)
			}

			err, ok := v.Errors["page_size"]
			if !ok && !tc.isValid {
				t.Fatal("expected page_size field to exist in validation error")
			}

			if err != tc.wantErr {
				t.Errorf("expected page_size error to be %s; got %s", tc.wantErr, err)
			}
		})
	}
}
