package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MangoLambda/KeyLox-Server/models"
	"github.com/go-chi/chi/v5"
)

// @Summary Gets a specific user
// @Description Gets a specific user
// @Tags user
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} models.UserResponse
// @Success 204 {object} string
// @Failure 500 {object} string
// @Router /user/{username} [get]
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := chi.URLParam(r, "username")

		var user models.DBUser
		rows, err := db.Query("SELECT id, username, clientSalt FROM users WHERE username = ?", username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.Scan(&user.ID, &user.Username, &user.ClientSalt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var userResponse models.UserResponse
			userResponse.Salt = user.ClientSalt
			json.NewEncoder(w).Encode(userResponse)
		} else {
			http.Error(w, "User not found", http.StatusNoContent)
		}
	}
}

// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User"
// @Success 200 {object} string
// @Failure 409 {object} string
// @Failure 500 {object} string
// @Router /register [post]
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the username already exists
		var existingUser models.UserResponse
		err := db.QueryRow("SELECT clientSalt FROM users WHERE username = ?", req.Username).Scan(&existingUser.Salt)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err == nil {
			// Username already exists
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		// Set server salt, vault id, and last login
		user := models.DBUser{
			Username:   req.Username,
			Key:        req.Key,
			ClientSalt: req.ClientSalt,
			ServerSalt: "some-generated-server-salt", // Replace with actual salt generation logic
			HashedKey:  "some-hashed-key",            // Replace with actual key hashing logic
			VaultId:    1,                            // Replace with actual vault ID logic if needed
			LastLogin:  time.Now(),
		}

		_, err = db.Exec("INSERT INTO users (username, key, clientSalt, serverSalt, hashedKey, vaultId, lastLogin) VALUES (?, ?, ?, ?, ?, ?, ?)",
			user.Username, user.Key, user.ClientSalt, user.ServerSalt, user.HashedKey, user.VaultId, user.LastLogin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("User registered successfully")
	}
}

// @Summary Get a vault by ID
// @Description Get a vault by ID
// @Tags vaults
// @Produce json
// @Param id path int true "Vault ID"
// @Success 200 {object} models.VaultResponse
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /vault/{id} [get]
func GetVault(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid vault ID", http.StatusBadRequest)
			return
		}

		var vault models.DBVault
		err = db.QueryRow("SELECT id, userId, filename, created_at FROM vaults WHERE id = ?", id).Scan(&vault.ID, &vault.UserId, &vault.FileName, &vault.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Vault not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(vault)
	}
}
