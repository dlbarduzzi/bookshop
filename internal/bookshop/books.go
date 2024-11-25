package bookshop

import (
	"fmt"
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsontil"
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

type listBooksResponse struct {
	Code  int           `json:"code"`
	Books []*model.Book `json:"books"`
}

func (b *Bookshop) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		model.Filters
	}

	q := r.URL.Query()
	v := validator.NewValidator()

	input.Filters.Page = b.readInt(q, "page", 1, v)
	input.Filters.PageSize = b.readInt(q, "page_size", 10, v)

	if !v.IsValid() {
		b.serverError(w, r, fmt.Errorf("invalid input"))
		return
	}

	books, err := b.models.Books.GetAll(input.Filters)
	if err != nil {
		b.serverError(w, r, err)
		return
	}

	res := listBooksResponse{
		Code:  http.StatusOK,
		Books: books,
	}

	if err := jsontil.Marshal(w, res, res.Code, nil); err != nil {
		b.serverError(w, r, err)
		return
	}
}
