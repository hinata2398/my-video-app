package persistence

import (
	"database/sql"
	"errors"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash, username string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash, username) VALUES ($1, $2, $3)
		 RETURNING id, email, username, password_hash, created_at`,
		email, passwordHash, username,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, errors.New("このメールアドレスはすでに登録されています")
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(
		`SELECT id, email, username, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
