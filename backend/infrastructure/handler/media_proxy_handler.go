package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hinata2398/my-video-app/backend/infrastructure/storage"
)

type MediaHandler struct {
	minio *storage.MinioClient
}

func NewMediaHandler(minio *storage.MinioClient) *MediaHandler {
	return &MediaHandler{minio: minio}
}

func (h *MediaHandler) Serve(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "*") // /media/ より後ろ全部
	obj, err := h.minio.GetObject(r.Context(), key)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	defer obj.Close()

	stat, err := obj.Stat() // 実取得 → 無ければエラー
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if stat.ContentType != "" {
		w.Header().Set("Content-Type", stat.ContentType)
	}
	http.ServeContent(w, r, key, stat.LastModified, obj) // Range自動処理
}

type ProxyResolver struct {
	baseURL string // 例: http://localhost:8080
}

func NewProxyResolver(baseURL string) *ProxyResolver {
	return &ProxyResolver{baseURL: baseURL}
}

func (p *ProxyResolver) PublicURL(key string) string {
	if key == "" {
		return ""
	}
	return p.baseURL + "/media/" + key
}

type CloudFrontResolver struct {
	baseURL string // 例: https://dxxxx.cloudfront.net
}

func NewCloudFrontResolver(baseURL string) *CloudFrontResolver {
	return &CloudFrontResolver{baseURL: baseURL}
}

func (c *CloudFrontResolver) PublicURL(key string) string {
	if key == "" {
		return ""
	}
	return c.baseURL + "/" + key
}
