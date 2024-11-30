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
// @Success 200 {object} models.User
// @Failure 500 {object} map[string]string
// @Router /user/{username} [get]
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := chi.URLParam(r, "username")

		var user models.User
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
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "User not found", http.StatusNotFound)
		}
	}
}

// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 200 {object} models.User
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		json.NewDecoder(r.Body).Decode(&user)
		// Check if the username already exists and get clientSalt if it does
		var clientSalt sql.NullString
		err := db.QueryRow("SELECT clientSalt FROM users WHERE username = ?", user.Username).Scan(&clientSalt)
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if clientSalt.Valid {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"clientSalt": clientSalt.String})
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		// Set server salt, vault id, and last login
		user.ServerSalt = "some-generated-server-salt" // Replace with actual salt generation logic
		user.VaultId = 0                               // Replace with actual vault ID logic if needed
		user.LastLogin = time.Now()
		result, err := db.Exec("INSERT INTO users (username, key, clientSalt, serverSalt, vaultId, lastLogin) VALUES (?, ?, ?, ?, ?, ?)",
			user.Username, user.Key, user.ClientSalt, user.ServerSalt, user.VaultId, user.LastLogin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.ID = int(id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// @Summary Get a vault by ID
// @Description Get a vault by ID
// @Tags vaults
// @Produce json
// @Param id path int true "Vault ID"
// @Success 200 {object} models.Vault
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
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

		var vault models.Vault
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
