package bookshop

import (
	"errors"
	"net/http"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/bookshop/model"
	"github.com/dlbarduzzi/bookshop/internal/jsoner"
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

func (bs *Bookshop) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := jsoner.Unmarshal(w, r, &input); err != nil {
		bs.clientError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user := &model.User{
		Name:          input.Name,
		Email:         input.Email,
		EmailVerified: false,
	}

	password := &model.Password{
		Plaintext: input.Password,
	}

	v := validator.NewValidator()

	user.Validate(v)
	password.Validate(v)

	if !v.IsValid() {
		bs.validationError(w, r, v.Errors)
		return
	}

	if err := password.Encrypt(); err != nil {
		bs.serverError(w, r, err)
		return
	}

	found, err := bs.models.Users.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, model.ErrRecordNotFound) {
		bs.serverError(w, r, err)
		return
	}

	if found != nil {
		bs.clientError(w, r, http.StatusBadRequest, "This email address is already in use.")
		return
	}

	user, err = bs.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			bs.clientError(w, r, http.StatusBadRequest, "This email address is already in use.")
		} else {
			bs.serverError(w, r, err)
		}
		return
	}

	if err := bs.models.Users.SavePassword(user.ID, password.Hash); err != nil {
		bs.serverError(w, r, err)
		return
	}

	_, err = bs.models.Tokens.Save(user.ID, time.Hour*24*1, model.ScopeEmailVerification)
	if err != nil {
		bs.serverError(w, r, err)
		return
	}

	data := jsoner.Envelope{
		"user": user,
	}

	if err := jsoner.Marshal(w, data, http.StatusCreated, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}

func (bs *Bookshop) verifyEmailUserHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("token")
	token := &model.Token{Plaintext: input}

	v := validator.NewValidator()

	if token.Validate(v); !v.IsValid() {
		bs.validationError(w, r, v.Errors)
		return
	}

	user, err := bs.models.Users.GetForToken(token.Plaintext, model.ScopeEmailVerification)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			bs.clientError(w, r, http.StatusBadRequest, "This token is invalid or have expired.")
		} else {
			bs.serverError(w, r, err)
		}
		return
	}

	user.EmailVerified = true

	err = bs.models.Users.Update(user)
	if err != nil {
		if errors.Is(err, model.ErrEditConflict) {
			bs.conflictError(w, r, "There was a conflict error when updating user! Please try again.")
		} else {
			bs.serverError(w, r, err)
		}
		return
	}

	err = bs.models.Tokens.DeleteAllForUser(user.ID, model.ScopeEmailVerification)
	if err != nil {
		bs.serverError(w, r, err)
		return
	}

	data := jsoner.Envelope{
		"user": user,
	}

	if err := jsoner.Marshal(w, data, http.StatusOK, nil); err != nil {
		bs.serverError(w, r, err)
		return
	}
}
