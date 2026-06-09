package entity

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	VideoID   int64     `json:"video_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Body      string    `json:"body"`
	LikeCount    int64     `json:"like_count"`
	Liked        bool      `json:"liked"`
	DislikeCount int64     `json:"dislike_count"`
	Disliked     bool      `json:"disliked"`
	CreatedAt    time.Time `json:"created_at"`
}
