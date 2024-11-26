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
		model.Filters
	}

	q := r.URL.Query()
	v := validator.NewValidator()

	input.Filters.Page = b.readInt(q, "page", 1, v)
	input.Filters.PageSize = b.readInt(q, "page_size", 10, v)

	if input.Filters.Validate(v); !v.IsValid() {
		b.validationError(w, r, v.Errors)
		return
	}

	books, metadata, err := b.models.Books.GetAll(input.Filters)
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
