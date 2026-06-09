package entity

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	VideoID   int64     `json:"video_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}
