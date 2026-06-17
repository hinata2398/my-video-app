package persistence

import (
	"database/sql"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
)

type WatchHistoryRepository struct {
	db *sql.DB
}

func NewWatchHistoryRepository(db *sql.DB) *WatchHistoryRepository {
	return &WatchHistoryRepository{db: db}
}

func (r *WatchHistoryRepository) Add(userID, videoID int64) error {
	_, err := r.db.Exec(
		`INSERT INTO watch_histories (user_id, video_id, watched_at)
         VALUES ($1, $2, NOW())
         ON CONFLICT (user_id, video_id) DO UPDATE SET watched_at = NOW()`,
		userID, videoID,
	)
	return err
}

func (r *WatchHistoryRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	rows, err := r.db.Query(`
		SELECT v.id, v.user_id, v.title, v.description, v.thumbnail_url, v.video_url,
		       v.status, v.status_message, v.view_count, v.created_at, v.updated_at
		FROM videos v
		JOIN watch_histories b ON b.video_id = v.id
		WHERE b.user_id = $1
		ORDER BY b.watched_at DESC`,
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
