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

	userRepo := persistence.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	minioClient, minioErr := storage.NewMinioClient()
	if minioErr != nil {
		log.Fatal("Could not connect to MinIO:", minioErr)
	}

	videoRepo := persistence.NewVideoRepository(db)
	videoUsecase := usecase.NewVideoUsecase(videoRepo)
	videoHandler := handler.NewVideoHandler(videoUsecase)
	uploadHandler := handler.NewUploadHandler(minioClient)
	thumbnailHandler := handler.NewThumbnailHandler(minioClient, db)
	transcodeQueue := queue.NewTranscodeQueue(minioClient, db, 2) // worker 2本
	transcodeHandler := handler.NewTranscodeHandler(transcodeQueue, db)

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

	// 作成・更新・削除・アップロードは認証必要
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Auth)
		r.Get("/api/me/videos", videoHandler.MyList)
		r.Post("/api/videos", videoHandler.Create)
		r.Put("/api/videos/{id}", videoHandler.Update)
		r.Delete("/api/videos/{id}", videoHandler.Delete)
		r.Get("/api/videos/{id}/upload-url", uploadHandler.PresignedURL)
		r.Get("/api/videos/{id}/thumbnail-upload-url", uploadHandler.PresignedThumbnailURL)
		r.Post("/api/videos/{id}/generate-thumbnail", thumbnailHandler.Generate)
		r.Post("/api/videos/{id}/transcode", transcodeHandler.Enqueue)
	})
	// ステータス確認は認証不要（一覧ページからも参照できるように）
	r.Get("/api/videos/{id}/status", transcodeHandler.Status)

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
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

		CREATE TABLE IF NOT EXISTS videos (
			id            BIGSERIAL PRIMARY KEY,
			user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title         TEXT NOT NULL,
			description   TEXT NOT NULL DEFAULT '',
			thumbnail_url TEXT NOT NULL DEFAULT '',
			video_url     TEXT NOT NULL DEFAULT '',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_videos_created_at ON videos (created_at DESC);
		CREATE INDEX IF NOT EXISTS idx_videos_user_id ON videos (user_id);
		ALTER TABLE videos ADD COLUMN IF NOT EXISTS video_url TEXT NOT NULL DEFAULT '';
		ALTER TABLE videos ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'done';
		ALTER TABLE videos ADD COLUMN IF NOT EXISTS status_message TEXT NOT NULL DEFAULT '';
	`)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration complete")
}
