package bookshop

import (
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsontil"
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

type listBooksResponse struct {
	Code     int            `json:"code"`
	Books    []*model.Book  `json:"books"`
	Metadata model.Metadata `json:"metadata"`
}

func (b *Bookshop) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string
		Categories []string
		model.Filters
	}

	q := r.URL.Query()
	v := validator.NewValidator()

	input.Title = b.readString(q, "title", "")
	input.Categories = b.readCSV(q, "categories", []string{})

	input.Filters.Page = b.readInt(q, "page", 1, v)
	input.Filters.PageSize = b.readInt(q, "page_size", 10, v)

	input.Filters.Sort = b.readString(q, "sort", "id")
	input.Filters.SortSafeList = []string{
		"id", "title", "published_date", "page_count",
		"-id", "-title", "-published_date", "-page_count",
	}

	if input.Filters.Validate(v); !v.IsValid() {
		b.validationError(w, r, v.Errors)
		return
	}

	books, metadata, err := b.models.Books.GetAll(input.Title, input.Categories, input.Filters)
	if err != nil {
		b.serverError(w, r, err)
		return
	}

	res := listBooksResponse{
		Code:     http.StatusOK,
		Books:    books,
		Metadata: metadata,
	}

	if err := jsontil.Marshal(w, res, res.Code, nil); err != nil {
		b.serverError(w, r, err)
		return
	}
}
