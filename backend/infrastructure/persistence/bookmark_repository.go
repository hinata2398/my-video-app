package persistence

import (
	"database/sql"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
)

type BookmarkRepository struct {
	db *sql.DB
}

func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Add(userID, videoID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO bookmarks (user_id, video_id) VALUES ($1, $2)`,
		userID, videoID,
	)
	return err
}

func (r *BookmarkRepository) Remove(userID, videoID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM bookmarks WHERE user_id = $1 AND video_id = $2`,
		userID, videoID,
	)
	return err
}

func (r *BookmarkRepository) Exists(userID, videoID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(
            SELECT 1 FROM bookmarks WHERE user_id = $1 AND video_id = $2
        )`,
		userID, videoID,
	).Scan(&exists)
	return exists, err
}

func (r *BookmarkRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	rows, err := r.db.Query(`
		SELECT v.id, v.user_id, v.title, v.description, v.thumbnail_url, v.video_url,
		       v.status, v.status_message, v.view_count, v.created_at, v.updated_at
		FROM videos v
		JOIN bookmarks b ON b.video_id = v.id
		WHERE b.user_id = $1
		ORDER BY b.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		v := &entity.Video{}
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.ThumbnailURL, &v.VideoURL,
			&v.Status, &v.StatusMessage, &v.ViewCount, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}
