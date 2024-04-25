package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/AdityaVarmaUddaraju/paytm/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("username already exists")
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  password  `json:"-"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)

	if err != nil {
		return err
	}

	p.plaintext = &plainText
	p.hash = hash

	return nil
}

func (p *password) Match(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}


func ValidatePassword(v *validator.Validator, plainTextPassword string) {
	v.ValidateEmpty(plainTextPassword, "password")
	v.Check(len(plainTextPassword) <= 72, "password", "password cannot be greater than 72 bytes")
}

func ValidateUser(v *validator.Validator, user *User) {
	
	v.ValidateEmpty(user.UserName, "username")
	v.ValidateEmpty(user.FirstName, "firstname")
	v.ValidateEmpty(user.LastName, "lastname")

	ValidatePassword(v, *user.Password.plaintext)

	if user.Password.hash == nil {
		panic("missing password hash for the user")
	}
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, firstname, lastname, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{user.UserName, user.FirstName, user.LastName, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"":
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) GetByUsername(firstName string) (*User, error) {
	query := `
		SELECT id, created_at, username, firstname, lastname, password_hash, version
		FROM users
		WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, firstName).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UserName,
		&user.FirstName,
		&user.LastName,
		&user.Password.hash,
		&user.Version,
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

func (m *UserModel) GetUsers(searchTerm string) ([]*User, error) {
	query := `
	SELECT id, created_at, username, firstname, lastname, password_hash, version
	FROM users
	WHERE username like $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)

	defer cancel()

	var users []*User

	rows, err := m.DB.QueryContext(ctx, query, fmt.Sprintf("%%%s%%", searchTerm))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User

		rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UserName,
			&user.FirstName,
			&user.LastName,
			&user.Password.hash,
			&user.Version,
		)

		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username = $1, firstname = $2, lastname = $3, password_hash = $4, version = version + 1
	WHERE id = $5 AND version = $6
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	args := []interface{}{user.UserName, user.FirstName, user.LastName, user.Password.hash, user.ID, user.Version}

	_, err  := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

