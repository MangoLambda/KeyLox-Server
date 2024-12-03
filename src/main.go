package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/MangoLambda/KeyLox-Server/src/docs"
	getHandlers "github.com/MangoLambda/KeyLox-Server/src/handlers/get"
	postHandlers "github.com/MangoLambda/KeyLox-Server/src/handlers/post"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title KeyLox Server API
// @version 1.0
// @description Credentials synchronization server for KeyLox.
// @license.name MIT
// @license.url https://mit-license.org/
// @BasePath /
// @securityDefinitions.basic BasicAuth
func main() {
	db := setupDb()
	defer db.Close()

	r := chi.NewRouter()
	setupMiddleware(r)
	setupRoutes(r, db)

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	log.Fatal(srv.ListenAndServe())
}

func setupDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}

	// Enable foreign key support
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		panic(err)
	}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		panic(err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		client_salt TEXT,
		server_salt TEXT,
		hashed_key TEXT,
		last_login DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vaults (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		panic(err)
	}

	// Periodically trigger a checkpoint
	go func() {
		for {
			_, err := db.Exec("PRAGMA wal_checkpoint(FULL);")
			if err != nil {
				log.Printf("Error during checkpoint: %v", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	return db
}

func setupMiddleware(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(httprate.LimitByIP(100, 1*time.Minute))
}

func setupRoutes(router *chi.Mux, db *sql.DB) {
	router.Post("/register", postHandlers.RegisterHandler(db))
	router.Post("/vault/{username}", postHandlers.VaultHandler(db))
	router.Get("/user/{username}", getHandlers.GetUserHandler(db))
	router.Get("/vault/{username}", getHandlers.GetVaultHandler(db))

	router.Get("/swagger/*", httpSwagger.WrapHandler)
}
