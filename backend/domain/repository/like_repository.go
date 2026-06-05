package repository

type LikeRepository interface {
    Like(userID, videoID int64) error    // いいねする
    Unlike(userID, videoID int64) error  // いいねを取り消す
    Count(videoID int64) (int64, error)  // 何人いいねしたか数える
    Exists(userID, videoID int64) (bool, error) // 自分がいいね済みか確認
}
