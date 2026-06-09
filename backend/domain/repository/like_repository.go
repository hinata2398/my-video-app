package repository

type LikeRepository interface {
    Like(userID, videoID int64) error
    Unlike(userID, videoID int64) error
    Count(videoID int64) (int64, error)
    Exists(userID, videoID int64) (bool, error)
    Dislike(userID, videoID int64) error
    Undislike(userID, videoID int64) error
    DislikeCount(videoID int64) (int64, error)
    DislikeExists(userID, videoID int64) (bool, error)
}
