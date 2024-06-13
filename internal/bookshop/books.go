package bookshop

import (
	"fmt"
	"net/http"
)

func (bs *Bookshop) createBookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new book")
}

func (bs *Bookshop) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bs.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show book with id %d\n", id)
}
