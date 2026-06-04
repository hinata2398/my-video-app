package persistence

import (
	"database/sql"
	"errors"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
)

type VideoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) Create(userID int64, title, description, thumbnailURL string) (*entity.Video, error) {
	video := &entity.Video{}
	err := r.db.QueryRow(
		`INSERT INTO videos (user_id, title, description, thumbnail_url)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, title, description, thumbnail_url, video_url, created_at, updated_at`,
		userID, title, description, thumbnailURL,
	).Scan(&video.ID, &video.UserID, &video.Title, &video.Description, &video.ThumbnailURL, &video.VideoURL, &video.CreatedAt, &video.UpdatedAt)
	return video, err
}

func (r *VideoRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, description, thumbnail_url, video_url, created_at, updated_at
		 FROM videos WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		v := &entity.Video{}
		if err := rows.Scan(&v.ID, &v.UserID, &v.Title, &v.Description, &v.ThumbnailURL, &v.VideoURL, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (r *VideoRepository) FindAll() ([]*entity.Video, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, description, thumbnail_url, video_url, created_at, updated_at
		 FROM videos ORDER BY created_at DESC LIMIT 50`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		v := &entity.Video{}
		if err := rows.Scan(&v.ID, &v.UserID, &v.Title, &v.Description, &v.ThumbnailURL, &v.VideoURL, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (r *VideoRepository) FindByID(id int64) (*entity.Video, error) {
	video := &entity.Video{}
	err := r.db.QueryRow(
		`SELECT id, user_id, title, description, thumbnail_url, video_url, created_at, updated_at
		 FROM videos WHERE id = $1`,
		id,
	).Scan(&video.ID, &video.UserID, &video.Title, &video.Description, &video.ThumbnailURL, &video.VideoURL, &video.CreatedAt, &video.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	return video, err
}

func (r *VideoRepository) Update(id, userID int64, title, description, thumbnailURL, videoURL string) (*entity.Video, error) {
	video := &entity.Video{}
	err := r.db.QueryRow(
		`UPDATE videos SET title=$1, description=$2, thumbnail_url=$3, video_url=$4, updated_at=NOW()
		 WHERE id=$5 AND user_id=$6
		 RETURNING id, user_id, title, description, thumbnail_url, video_url, created_at, updated_at`,
		title, description, thumbnailURL, videoURL, id, userID,
	).Scan(&video.ID, &video.UserID, &video.Title, &video.Description, &video.ThumbnailURL, &video.VideoURL, &video.CreatedAt, &video.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("not found or forbidden")
	}
	return video, err
}

func (r *VideoRepository) UpdateVideoURL(id int64, videoURL string) error {
	_, err := r.db.Exec(
		`UPDATE videos SET video_url=$1, updated_at=NOW() WHERE id=$2`,
		videoURL, id,
	)
	return err
}

func (r *VideoRepository) Delete(id, userID int64) error {
	result, err := r.db.Exec(
		`DELETE FROM videos WHERE id=$1 AND user_id=$2`, id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return errors.New("not found or forbidden")
	}
	return nil
}
