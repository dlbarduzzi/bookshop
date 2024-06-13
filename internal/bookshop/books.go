package bookshop

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsoner"
)

func (bs *Bookshop) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string          `json:"title"`
		Authors       []string        `json:"authors"`
		PublishedDate string          `json:"published_date"`
		PageCount     model.PageCount `json:"page_count"`
		Categories    []string        `json:"categories"`
	}

	code, err := jsoner.Unmarshal(w, r, &input)
	if err != nil {
		if code == http.StatusInternalServerError {
			bs.serverError(w, r, err)
		} else {
			bs.clientError(w, r, code, err.Error())
		}
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
