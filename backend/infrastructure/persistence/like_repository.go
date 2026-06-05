package persistence

import (
	"database/sql"
)

type LikeRepository struct {
	db *sql.DB
}

func NewLikeRepository(db *sql.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Like(userID, videoID int64) error {
    _, err := r.db.Exec(
        `INSERT INTO likes (user_id, video_id) VALUES ($1, $2)`,
        userID, videoID,
    )
    return err
}

func (r *LikeRepository) Unlike(userID, videoID int64) error {
    _, err := r.db.Exec(
        `DELETE FROM likes WHERE user_id = $1 AND video_id = $2`,
        userID, videoID,
    )
    return err
}

func (r *LikeRepository) Count(videoID int64) (int64, error) {
    var count int64
    err := r.db.QueryRow(
        `SELECT COUNT(*) FROM likes WHERE video_id = $1`,
        videoID,
    ).Scan(&count)
    return count, err
}

func (r *LikeRepository) Exists(userID, videoID int64) (bool, error) {
    var exists bool
    err := r.db.QueryRow(
        `SELECT EXISTS(
            SELECT 1 FROM likes WHERE user_id = $1 AND video_id = $2
        )`,
        userID, videoID,
    ).Scan(&exists)
    return exists, err
}
