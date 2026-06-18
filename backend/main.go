package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hinata2398/my-video-app/backend/infrastructure/handler"
	authMiddleware "github.com/hinata2398/my-video-app/backend/infrastructure/middleware"
	"github.com/hinata2398/my-video-app/backend/infrastructure/persistence"
	"github.com/hinata2398/my-video-app/backend/infrastructure/queue"
	"github.com/hinata2398/my-video-app/backend/infrastructure/storage"
	"github.com/hinata2398/my-video-app/backend/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := connectDB()
	defer db.Close()

	minioClient, minioErr := storage.NewMinioClient()
	if minioErr != nil {
		log.Fatal("Could not connect to MinIO:", minioErr)
	}

	mediaHandler := handler.NewMediaHandler(minioClient)
	mediaResolver := handler.NewProxyResolver(os.Getenv("BACKEND_PUBLIC_URL"))

	userRepo := persistence.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase, mediaResolver)

	videoRepo := persistence.NewVideoRepository(db)
	videoUsecase := usecase.NewVideoUsecase(videoRepo)
	videoHandler := handler.NewVideoHandler(videoUsecase, mediaResolver)
	uploadHandler := handler.NewUploadHandler(minioClient)
	thumbnailHandler := handler.NewThumbnailHandler(minioClient, db, mediaResolver)
	transcodeQueue := queue.NewTranscodeQueue(minioClient, db, 2) // worker 2本
	transcodeHandler := handler.NewTranscodeHandler(transcodeQueue, db)
	likeRepo := persistence.NewLikeRepository(db)
	likeUsecase := usecase.NewLikeUsecase(likeRepo)
	likeHandler := handler.NewLikeHandler(likeUsecase)
	bookmarkRepo := persistence.NewBookmarkRepository(db)
	bookmarkUsecase := usecase.NewBookmarkUsecase(bookmarkRepo)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkUsecase, mediaResolver)

	commentRepo := persistence.NewCommentRepository(db)
	commentUsecase := usecase.NewCommentUsecase(commentRepo)
	commentHandler := handler.NewCommentHandler(commentUsecase)
	watchHistoryRepo := persistence.NewWatchHistoryRepository(db)
	watchHistoryUsecase := usecase.NewWatchHistoryUsecase(watchHistoryRepo)
	watchHistoryHandler := handler.NewWatchHistoryHandler(watchHistoryUsecase, mediaResolver)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/api/health", handler.Health)
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// 動画一覧・詳細は認証不要
	r.Get("/api/videos", videoHandler.List)
	r.Get("/api/videos/{id}", videoHandler.Get)
	r.Get("/media/*", mediaHandler.Serve)

	// 作成・更新・削除・アップロードは認証必要
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Auth)
		r.Get("/api/me", userHandler.GetMe)
		r.Put("/api/me", userHandler.UpdateMe)
		r.Get("/api/me/avatar-upload-url", uploadHandler.PresignedAvatarURL)
		r.Get("/api/me/videos", videoHandler.MyList)
		r.Post("/api/videos", videoHandler.Create)
		r.Put("/api/videos/{id}", videoHandler.Update)
		r.Delete("/api/videos/{id}", videoHandler.Delete)
		r.Get("/api/videos/{id}/upload-url", uploadHandler.PresignedURL)
		r.Get("/api/videos/{id}/thumbnail-upload-url", uploadHandler.PresignedThumbnailURL)
		r.Post("/api/videos/{id}/generate-thumbnail", thumbnailHandler.Generate)
		r.Post("/api/videos/{id}/transcode", transcodeHandler.Enqueue)
		r.Post("/api/videos/{id}/like", likeHandler.Toggle)
		r.Post("/api/videos/{id}/dislike", likeHandler.ToggleDislike)
		r.Get("/api/videos/{id}/like-status", likeHandler.Status)
	})
	// ステータス確認は認証不要（一覧ページからも参照できるように）
	r.Get("/api/videos/{id}/status", transcodeHandler.Status)

	// いいね数・よくないね数の取得は誰でも見られる（認証不要）
	r.Get("/api/videos/{id}/likes", likeHandler.Count)
	r.Get("/api/videos/{id}/dislikes", likeHandler.DislikeCount)

	r.Post("/api/videos/{id}/view", videoHandler.IncrementViewCount)

	// コメント（一覧は認証不要、投稿・いいねは認証必要）
	r.Get("/api/videos/{id}/comments", commentHandler.List)
	r.With(authMiddleware.Auth).Post("/api/videos/{id}/comments", commentHandler.Create)
	r.With(authMiddleware.Auth).Post("/api/comments/{commentId}/like", commentHandler.ToggleLike)
	r.With(authMiddleware.Auth).Post("/api/comments/{commentId}/dislike", commentHandler.ToggleDislike)

	// ブックマークの取得は認証必要
	r.With(authMiddleware.Auth).Post("/api/videos/{id}/bookmark", bookmarkHandler.Toggle)
	r.With(authMiddleware.Auth).Get("/api/videos/{id}/bookmark-status", bookmarkHandler.Status)
	r.With(authMiddleware.Auth).Get("/api/bookmarks", bookmarkHandler.FindByUserID)

	// 視聴履歴の取得は認証必要
	r.With(authMiddleware.Auth).Post("/api/videos/{id}/watch-history", watchHistoryHandler.Add)
	r.With(authMiddleware.Auth).Get("/api/watch-history", watchHistoryHandler.FindByUserID)

	log.Println("Backend running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func connectDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Println("Connected to database")
				migrate(db)
				return db
			}
		}
		log.Printf("Waiting for database... (%d/10)\n", i+1)
		time.Sleep(2 * time.Second)
	}
	log.Fatal("Could not connect to database:", err)
	return nil
}

func migrate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id            BIGSERIAL PRIMARY KEY,
			email         TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			username      TEXT NOT NULL DEFAULT '',
			avatar_url    TEXT NOT NULL DEFAULT '',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		ALTER TABLE users ADD COLUMN IF NOT EXISTS username   TEXT NOT NULL DEFAULT '';
		ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url TEXT NOT NULL DEFAULT '';
		CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

		CREATE TABLE IF NOT EXISTS videos (
			id             BIGSERIAL PRIMARY KEY,
			user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title          TEXT NOT NULL,
			description    TEXT NOT NULL DEFAULT '',
			thumbnail_url  TEXT NOT NULL DEFAULT '',
			video_url      TEXT NOT NULL DEFAULT '',
			status         TEXT NOT NULL DEFAULT 'done',
			status_message TEXT NOT NULL DEFAULT '',
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos (created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos (user_id);

		CREATE TABLE IF NOT EXISTS likes (
			id         BIGSERIAL PRIMARY KEY,
			user_id    BIGINT NOT NULL,
			video_id   BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS video_dislikes (
			id         BIGSERIAL PRIMARY KEY,
			user_id    BIGINT NOT NULL,
			video_id   BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, video_id)
		);

		CREATE TABLE IF NOT EXISTS comments (
			id         BIGSERIAL PRIMARY KEY,
			video_id   BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
			user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			body       TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_comments_video_id ON comments (video_id);

		CREATE TABLE IF NOT EXISTS comment_likes (
			id         BIGSERIAL PRIMARY KEY,
			comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
			user_id    BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(comment_id, user_id)
		);

		CREATE TABLE IF NOT EXISTS comment_dislikes (
			id         BIGSERIAL PRIMARY KEY,
			comment_id BIGINT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
			user_id    BIGINT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(comment_id, user_id)
		);

		CREATE TABLE IF NOT EXISTS bookmarks (
			id         BIGSERIAL PRIMARY KEY,
			user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			video_id   BIGINT NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, video_id)
		);

		CREATE TABLE IF NOT EXISTS watch_histories (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL REFERENCES users(id),
			video_id BIGINT NOT NULL REFERENCES videos(id),
			watched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			UNIQUE(user_id, video_id)
		);
	`)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration complete")
}
