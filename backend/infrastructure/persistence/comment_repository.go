package persistence

import (
	"database/sql"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(videoID, userID int64, body string) (*entity.Comment, error) {
	c := &entity.Comment{}
	err := r.db.QueryRow(`
		INSERT INTO comments (video_id, user_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, video_id, user_id, body, created_at`,
		videoID, userID, body,
	).Scan(&c.ID, &c.VideoID, &c.UserID, &c.Body, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	// メールアドレスを取得
	r.db.QueryRow(`SELECT email FROM users WHERE id = $1`, userID).Scan(&c.Email)
	return c, nil
}

func (r *CommentRepository) FindByVideoID(videoID int64) ([]*entity.Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.video_id, c.user_id, u.email, c.body, c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.video_id = $1
		ORDER BY c.created_at ASC`,
		videoID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		c := &entity.Comment{}
		if err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Email, &c.Body, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
