package bookshop

import (
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsontil"
)

type listBooksResponse struct {
	Code  int           `json:"code"`
	Books []*model.Book `json:"books"`
}

func (b *Bookshop) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		model.Filters
	}

	input.Filters.Page = 1
	input.Filters.PageSize = 10

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
