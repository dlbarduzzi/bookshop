package bookshop

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsoner"
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

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

	book := model.Book{
		ID:            id,
		Title:         "Skills Learning",
		Authors:       []string{"John Cooper"},
		PublishedDate: time.Now().Format("2006-01-02"),
		PageCount:     296,
		Categories:    []string{"software development", "improvement"},
		Version:       1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	data := jsoner.Envelope{
		"book": book,
	}

	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}
