package bookshop

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsoner"
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func (bs *Bookshop) listBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title      string
		Categories []string
		model.Filters
	}

	q := r.URL.Query()
	v := validator.NewValidator()

	input.Title = bs.readString(q, "title", "")
	input.Categories = bs.readCSV(q, "categories", []string{})

	input.Filters.Page = bs.readInt(q, "page", 1, v)
	input.Filters.PageSize = bs.readInt(q, "page_size", 20, v)

	input.Filters.Sort = bs.readString(q, "sort", "id")
	input.Filters.SortSafeList = []string{
		"id", "title", "published_date", "page_count",
		"-id", "-title", "-published_date", "-page_count",
	}

	if input.Filters.Validate(v); !v.IsValid() {
		bs.validationError(w, r, v.Errors)
		return
	}

	books, err := bs.models.Books.GetAll(input.Title, input.Categories, input.Filters)
	if err != nil {
		bs.serverError(w, r, err)
		return
	}

	data := jsoner.Envelope{
		"books": books,
	}

	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}

func (bs *Bookshop) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string          `json:"title"`
		Authors       []string        `json:"authors"`
		PublishedDate string          `json:"published_date"`
		PageCount     model.PageCount `json:"page_count"`
		Categories    []string        `json:"categories"`
	}

	if err := jsoner.Unmarshal(w, r, &input); err != nil {
		bs.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	book := &model.Book{
		Title:         input.Title,
		Authors:       input.Authors,
		PublishedDate: input.PublishedDate,
		PageCount:     input.PageCount,
		Categories:    input.Categories,
	}

	v := validator.NewValidator()

	if book.Validate(v); !v.IsValid() {
		bs.validationError(w, r, v.Errors)
		return
	}

	if err := bs.models.Books.Insert(book); err != nil {
		bs.serverError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/books/%d", book.ID))

	data := jsoner.Envelope{
		"book": book,
	}

	if err := jsoner.Marshal(w, data, http.StatusCreated, headers); err != nil {
		bs.serverError(w, r, err)
		return
	}
}

func (bs *Bookshop) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bs.readIDParam(r)
	if err != nil {
		bs.clientError(w, r, http.StatusNotFound, "Book with given id was not found.")
		return
	}

	book, err := bs.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			bs.clientError(w, r, http.StatusNotFound, "Book with given id was not found.")
		default:
			bs.serverError(w, r, err)
		}
		return
	}

	data := jsoner.Envelope{
		"book": book,
	}

	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}

func (bs *Bookshop) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bs.readIDParam(r)
	if err != nil {
		bs.clientError(w, r, http.StatusNotFound, "Book with given id was not found.")
		return
	}

	book, err := bs.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			bs.clientError(w, r, http.StatusNotFound, "Book with given id was not found.")
		default:
			bs.serverError(w, r, err)
		}
		return
	}

	var input struct {
		Title         *string          `json:"title"`
		Authors       []string         `json:"authors"`
		PublishedDate *string          `json:"published_date"`
		PageCount     *model.PageCount `json:"page_count"`
		Categories    []string         `json:"categories"`
	}

	if err := jsoner.Unmarshal(w, r, &input); err != nil {
		bs.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Authors != nil {
		book.Authors = input.Authors
	}
	if input.PublishedDate != nil {
		book.PublishedDate = *input.PublishedDate
	}
	if input.PageCount != nil {
		book.PageCount = *input.PageCount
	}
	if input.Categories != nil {
		book.Categories = input.Categories
	}

	v := validator.NewValidator()

	if book.Validate(v); !v.IsValid() {
		bs.validationError(w, r, v.Errors)
		return
	}

	if err := bs.models.Books.Update(book); err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			bs.clientError(w, r, http.StatusConflict, "Book was not updated due to an edit conflict.")
		default:
			bs.serverError(w, r, err)
		}
		return
	}

	data := jsoner.Envelope{
		"book": book,
	}

	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}

func (bs *Bookshop) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bs.readIDParam(r)
	if err != nil {
		bs.clientError(w, r, http.StatusNotFound, "Book with given id was not found.")
		return
	}

	err = bs.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			bs.clientError(w, r, http.StatusNotFound, "Book with given id was not found.")
		default:
			bs.serverError(w, r, err)
		}
		return
	}

	data := jsoner.Envelope{
		"message": "Book successfully deleted.",
	}

	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}
