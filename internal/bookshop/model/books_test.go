package model

import (
	"testing"

	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func TestBookValidate(t *testing.T) {
	t.Parallel()
	v := validator.NewValidator()

	book := &Book{
		Title:         "Test Title",
		Authors:       []string{"Test Author"},
		PublishedDate: "2024-01-30",
		PageCount:     100,
		Categories:    []string{"test"},
	}

	book.Validate(v)

	if !v.IsValid() {
		t.Error("expected book validation to succeed")
	}
}

func TestBookValidateTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		title   string
		isValid bool
		wantErr string
	}{
		{
			name:    "empty",
			title:   "",
			isValid: false,
			wantErr: "Title is required.",
		},
		{
			name:    "short",
			title:   "T",
			isValid: false,
			wantErr: "Title must have a value between 2 and 300 characters.",
		},
		{
			name:    "valid",
			title:   "Test",
			isValid: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Book{
				Title: tt.title,
			}

			v := validator.NewValidator()
			b.validateTitle(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected book title validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["title"]
			if !exists && !tt.isValid {
				t.Error("expected title field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected title error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}

func TestBookValidateAuthors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		authors []string
		isValid bool
		wantErr string
	}{
		{
			name:    "short",
			authors: []string{"T"},
			isValid: false,
			wantErr: "Authors names must have a value between 2 and 200 characters.",
		},
		{
			name:    "empty",
			authors: []string{},
			isValid: false,
			wantErr: "At least 1 author is required.",
		},
		{
			name:    "duplicate",
			authors: []string{"Test", "Test"},
			isValid: false,
			wantErr: "Authors cannot have duplicated values.",
		},
		{
			name:    "valid",
			authors: []string{"Test"},
			isValid: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Book{
				Authors: tt.authors,
			}

			v := validator.NewValidator()
			b.validateAuthors(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected book authors validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["authors"]
			if !exists && !tt.isValid {
				t.Error("expected authors field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected title error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}

func TestBookValidatePublishedDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		publishedDate string
		isValid       bool
		wantErr       string
	}{
		{
			name:          "empty",
			publishedDate: "",
			isValid:       false,
			wantErr:       "Published date is required.",
		},
		{
			name:          "invalid",
			publishedDate: "invalid",
			isValid:       false,
			wantErr:       "Invalid published date format (i.e. 2021-02-14).",
		},
		{
			name:          "future",
			publishedDate: "9999-01-30",
			isValid:       false,
			wantErr:       "Published date cannot be in the future.",
		},
		{
			name:          "valid",
			publishedDate: "2024-01-30",
			isValid:       true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Book{
				PublishedDate: tt.publishedDate,
			}

			v := validator.NewValidator()
			b.validatePublishedDate(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected book published date validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["published_date"]
			if !exists && !tt.isValid {
				t.Error("expected published_date field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected published_date error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}

func TestBookValidatePageCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pageCount PageCount
		isValid   bool
		wantErr   string
	}{
		{
			name:      "negative",
			pageCount: -1,
			isValid:   false,
			wantErr:   "Page count must have a value higher than 0.",
		},
		{
			name:      "zero",
			pageCount: 0,
			isValid:   false,
			wantErr:   "Page count must have a value higher than 0.",
		},
		{
			name:      "valid",
			pageCount: 100,
			isValid:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Book{
				PageCount: tt.pageCount,
			}

			v := validator.NewValidator()
			b.validatePageCount(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected book page count validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["page_count"]
			if !exists && !tt.isValid {
				t.Error("expected page_count field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected page_count error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}

func TestBookValidateCategories(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		categories []string
		isValid    bool
		wantErr    string
	}{
		{
			name:       "short",
			categories: []string{"t"},
			isValid:    false,
			wantErr:    "Categories must have values between 2 and 200 characters.",
		},
		{
			name:       "empty",
			categories: []string{},
			isValid:    false,
			wantErr:    "At least 1 category is required.",
		},
		{
			name:       "too_many",
			categories: []string{"t1", "t2", "t3", "t4", "t5", "t6"},
			isValid:    false,
			wantErr:    "Categories cannot have more than 5 values.",
		},
		{
			name:       "duplicate",
			categories: []string{"t1", "t1"},
			isValid:    false,
			wantErr:    "Categories cannot have duplicated values.",
		},
		{
			name:       "valid",
			categories: []string{"t1"},
			isValid:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Book{
				Categories: tt.categories,
			}

			v := validator.NewValidator()
			b.validateCategories(v)

			if v.IsValid() != tt.isValid {
				t.Errorf("expected book categories validation to be %t; got %t", tt.isValid, v.IsValid())
			}

			err, exists := v.Errors["categories"]
			if !exists && !tt.isValid {
				t.Error("expected categories field in book validation error")
			}

			if err != tt.wantErr {
				t.Errorf("expected categories error to be %s; got %s", tt.wantErr, err)
			}
		})
	}
}
