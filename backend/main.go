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
	"github.com/hinata2398/my-video-app/backend/infrastructure/persistence"
	"github.com/hinata2398/my-video-app/backend/usecase"
	_ "github.com/lib/pq"
)

func main() {
	db := connectDB()
	defer db.Close()

	userRepo := persistence.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/health", handler.Health)
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	log.Println("Backend running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
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
	`)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration complete")
}
