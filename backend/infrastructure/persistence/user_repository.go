package persistence

import (
	"database/sql"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(
		`INSERT INTO users (email, password_hash) VALUES ($1, $2)
		 RETURNING id, email, password_hash, created_at`,
		email, passwordHash,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRow(
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
