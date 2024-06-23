package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"strings"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/security"
	"github.com/dlbarduzzi/bookshop/internal/validator"
)

const ScopeEmailVerification = "email-verification"

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

type TokenStore interface {
	Save(userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllForUser(userID int64, scope string) error
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	plaintext, err := security.RandomString(54)
	if err != nil {
		return nil, err
	}

	token.Plaintext = plaintext

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func (m TokenModel) Save(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m TokenModel) DeleteAllForUser(userID int64, scope string) error {
	query := `
		DELETE FROM tokens 
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

func (t *Token) Validate(v *validator.Validator) {
	plaintext := strings.TrimSpace(t.Plaintext)

	if plaintext == "" {
		v.AddError("token", "Token is required.")
		return
	}

	if len(plaintext) != 54 {
		v.AddError("token", "Token must be 54 characters long.")
		return
	}

	// Re-assign clean value.
	t.Plaintext = plaintext
}
