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
	// ユーザー名を取得
	r.db.QueryRow(`SELECT username FROM users WHERE id = $1`, userID).Scan(&c.Username)
	return c, nil
}

func (r *CommentRepository) FindByVideoID(videoID, userID int64) ([]*entity.Comment, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.video_id, c.user_id, u.username, c.body,
		       COUNT(DISTINCT cl.id) AS like_count,
		       COALESCE(BOOL_OR(cl.user_id = $2), false) AS liked,
		       COUNT(DISTINCT cd.id) AS dislike_count,
		       COALESCE(BOOL_OR(cd.user_id = $2), false) AS disliked,
		       c.created_at
		FROM comments c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN comment_likes cl ON cl.comment_id = c.id
		LEFT JOIN comment_dislikes cd ON cd.comment_id = c.id
		WHERE c.video_id = $1
		GROUP BY c.id, u.username
		ORDER BY c.created_at ASC`,
		videoID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entity.Comment
	for rows.Next() {
		c := &entity.Comment{}
		if err := rows.Scan(&c.ID, &c.VideoID, &c.UserID, &c.Username, &c.Body, &c.LikeCount, &c.Liked, &c.DislikeCount, &c.Disliked, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepository) ToggleDislike(commentID, userID int64) (int64, bool, error) {
	var exists bool
	r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM comment_dislikes WHERE comment_id = $1 AND user_id = $2)`,
		commentID, userID,
	).Scan(&exists)

	if exists {
		r.db.Exec(`DELETE FROM comment_dislikes WHERE comment_id = $1 AND user_id = $2`, commentID, userID)
	} else {
		// dislike追加時はlikeを解除
		r.db.Exec(`DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2`, commentID, userID)
		r.db.Exec(`INSERT INTO comment_dislikes (comment_id, user_id) VALUES ($1, $2)`, commentID, userID)
	}

	var count int64
	r.db.QueryRow(`SELECT COUNT(*) FROM comment_dislikes WHERE comment_id = $1`, commentID).Scan(&count)

	return count, !exists, nil
}

func (r *CommentRepository) ToggleLike(commentID, userID int64) (int64, bool, error) {
	// すでにいいね済みか確認
	var exists bool
	r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = $1 AND user_id = $2)`,
		commentID, userID,
	).Scan(&exists)

	if exists {
		// いいね解除
		r.db.Exec(`DELETE FROM comment_likes WHERE comment_id = $1 AND user_id = $2`, commentID, userID)
	} else {
		// like追加時はdislikeを解除
		r.db.Exec(`DELETE FROM comment_dislikes WHERE comment_id = $1 AND user_id = $2`, commentID, userID)
		r.db.Exec(`INSERT INTO comment_likes (comment_id, user_id) VALUES ($1, $2)`, commentID, userID)
	}

	// 最新のいいね数を取得
	var count int64
	r.db.QueryRow(`SELECT COUNT(*) FROM comment_likes WHERE comment_id = $1`, commentID).Scan(&count)

	return count, !exists, nil
}
