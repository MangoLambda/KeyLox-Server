package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/MangoLambda/KeyLox-Server/docs"
	handlers "github.com/MangoLambda/KeyLox-Server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title KeyLox Server API
// @version 1.0
// @description Credentials synchronization server for KeyLox.

// @license.name MIT
// @license.url https://mit-license.org/

// @BasePath /

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create tables if they don't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE,
		key TEXT,
		clientSalt TEXT,
		serverSalt TEXT,
		hashedKey TEXT,
		vaultId INTEGER,
		lastLogin DATETIME
	)`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vaults (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userId INTEGER,
		filename TEXT,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/user/{username}", handlers.GetUser(db))
	r.Post("/register", handlers.Register(db))
	r.Get("/vault/{id}", handlers.GetVault(db))

	// Serve OpenAPI documentation
	r.Get("/swagger/*", httpSwagger.WrapHandler)

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
