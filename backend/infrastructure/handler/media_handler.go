package handler

import "github.com/hinata2398/my-video-app/backend/domain/entity"

// MediaURLResolver はオブジェクトキーを配信用URLに変換する。
type MediaURLResolver interface {
	PublicURL(key string) string
}

// resolveUser は ユーザーアイコンのURLを配信URLに変換する。
func resolveUser(u *entity.User, r MediaURLResolver) {
	if u == nil {
		return
	}
	u.AvatarURL = r.PublicURL(u.AvatarURL)
}

// resolveVideo は Video のキー項目(VideoURL/ThumbnailURL)を配信URLに変換する（破壊的）。
func resolveVideo(v *entity.Video, r MediaURLResolver) {
	if v == nil {
		return
	}
	v.VideoURL = r.PublicURL(v.VideoURL)
	v.ThumbnailURL = r.PublicURL(v.ThumbnailURL)
}

// resolveVideos はスライス内の全 Video を変換する。
func resolveVideos(vs []*entity.Video, r MediaURLResolver) {
	for _, v := range vs {
		resolveVideo(v, r)
	}
}
