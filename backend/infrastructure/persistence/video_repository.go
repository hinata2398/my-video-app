package persistence

import (
	"database/sql"
	"errors"

	"github.com/hinata2398/my-video-app/backend/domain/entity"
	"github.com/lib/pq"
)

type VideoRepository struct {
	db *sql.DB
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func scanVideo(row interface{ Scan(...interface{}) error }) (*entity.Video, error) {
	v := &entity.Video{}
	err := row.Scan(
		&v.ID, &v.UserID, &v.Title, &v.Description,
		&v.ThumbnailURL, &v.VideoURL, &v.Status, &v.StatusMessage,
		&v.CreatedAt, &v.UpdatedAt,
	)
	return v, err
}

const selectFields = `id, user_id, title, description, thumbnail_url, video_url, status, status_message, created_at, updated_at`

func (r *VideoRepository) Create(userID int64, title, description, thumbnailURL string) (*entity.Video, error) {
	row := r.db.QueryRow(
		`INSERT INTO videos (user_id, title, description, thumbnail_url)
		 VALUES ($1, $2, $3, $4)
		 RETURNING `+selectFields,
		userID, title, description, thumbnailURL,
	)
	return scanVideo(row)
}

func (r *VideoRepository) FindByUserID(userID int64) ([]*entity.Video, error) {
	rows, err := r.db.Query(
		`SELECT `+selectFields+` FROM videos WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		v, err := scanVideo(rows)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (r *VideoRepository) FindAll(query string) ([]*entity.Video, error) {
	var rows *sql.Rows
	var err error
	if query == "" {
		rows, err = r.db.Query(
			`SELECT `+selectFields+` FROM videos WHERE status = 'done' ORDER BY created_at DESC LIMIT 50`,
		)
	} else {
		rows, err = r.db.Query(
			`SELECT `+selectFields+` FROM videos
			 WHERE status = 'done'
			   AND (title ILIKE $1 OR description ILIKE $1)
			 ORDER BY created_at DESC LIMIT 50`,
			"%"+query+"%",
		)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []*entity.Video
	for rows.Next() {
		v, err := scanVideo(rows)
		if err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func (r *VideoRepository) FindByID(id int64) (*entity.Video, error) {
	row := r.db.QueryRow(
		`SELECT `+selectFields+` FROM videos WHERE id = $1`, id,
	)
	v, err := scanVideo(row)
	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	return v, err
}

func (r *VideoRepository) Update(id, userID int64, title, description, thumbnailURL, videoURL string) (*entity.Video, error) {
	row := r.db.QueryRow(
		`UPDATE videos SET title=$1, description=$2, thumbnail_url=$3, video_url=$4, updated_at=NOW()
		 WHERE id=$5 AND user_id=$6
		 RETURNING `+selectFields,
		title, description, thumbnailURL, videoURL, id, userID,
	)
	v, err := scanVideo(row)
	if err == sql.ErrNoRows {
		return nil, errors.New("not found or forbidden")
	}
	return v, err
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

var _ *pq.Error // pqを使用
