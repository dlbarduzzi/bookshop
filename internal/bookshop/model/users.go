package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/dlbarduzzi/bookshop/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	Version       int       `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

type UserStore interface {
	Insert(user *User) (*User, error)
	Update(user *User) error
	GetByEmail(email string) (*User, error)
	GetForToken(token string, scope string) (*User, error)
	SavePassword(userID int64, passwordHash []byte) error
}

func (m UserModel) Insert(user *User) (*User, error) {
	query := `
		INSERT INTO users (name, email, email_verified)
		VALUES ($1, $2, $3)
		RETURNING id, name, email, email_verified, version, created_at, updated_at`

	args := []any{user.Name, user.Email, user.EmailVerified}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, email_verified = $3, version = version + 1, updated_at = $4
        WHERE id = $5 AND version = $6
        RETURNING version, updated_at`

	args := []any{
		user.Name,
		user.Email,
		user.EmailVerified,
		time.Now().UTC(),
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version, &user.UpdatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, name, email, email_verified, version, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) SavePassword(userID int64, passwordHash []byte) error {
	query := `
		INSERT INTO passwords (hash, user_id)
		VALUES ($1, $2)`

	args := []any{passwordHash, userID}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m UserModel) GetForToken(token string, scope string) (*User, error) {
	hash := sha256.Sum256([]byte(token))

	query := `
		SELECT users.id, users.name, users.email, users.email_verified, users.version, users.created_at, users.updated_at
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3`

	args := []any{hash[:], scope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.EmailVerified,
		&user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *User) Validate(v *validator.Validator) {
	u.validateName(v)
	u.validateEmail(v)
}

func (u *User) validateName(v *validator.Validator) {
	name := strings.TrimSpace(u.Name)

	if name == "" {
		v.AddError("name", "Name is required.")
		return
	}

	if !validator.MatchRegex(name, regexp.MustCompile(`^[a-zA-Z\s]*$`)) {
		v.AddError("name", "Name must contain only letters and spaces.")
		return
	}

	if len(name) < 3 || len(name) > 200 {
		v.AddError("name", "Name must have a value between 3 and 200 characters.")
		return
	}

	// Re-assign clean value.
	u.Name = name
}

func (u *User) validateEmail(v *validator.Validator) {
	email := strings.TrimSpace(u.Email)

	if email == "" {
		v.AddError("email", "Email is required.")
		return
	}

	if !validator.MatchRegex(email, validator.EmailRegex) {
		v.AddError("email", "Invalid email.")
		return
	}

	// Re-assign clean value.
	u.Email = email
}

type Password struct {
	Plaintext string
	Hash      []byte
	UserID    int64
}

func (p *Password) Encrypt() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.Plaintext), 12)
	if err != nil {
		return err
	}
	p.Hash = hash
	return nil
}

func (p *Password) Validate(v *validator.Validator) {
	plaintext := strings.TrimSpace(p.Plaintext)

	if plaintext == "" {
		v.AddError("password", "Password is required.")
		return
	}

	if len(plaintext) < 8 || len(plaintext) > 72 {
		v.AddError("password", "Password must have a value between 8 and 72 characters.")
		return
	}

	// Re-assign clean value.
	p.Plaintext = plaintext
}
