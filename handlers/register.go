package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MangoLambda/KeyLox-Server/models"
)

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
		var registerRequest models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if the username already exists
		var exists, err = checkUsernameExists(db, registerRequest.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if exists {
			// Username already exists
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		// Register the user
		err = registerUser(db, registerRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("User registered successfully")
	}
}

func checkUsernameExists(db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func registerUser(db *sql.DB, registerRequest models.RegisterRequest) error {
	// Set server salt, vault id, and last login
	user := models.DBUser{
		Username:   registerRequest.Username,
		ClientSalt: registerRequest.ClientSalt,
		ServerSalt: "some-generated-server-salt", // Replace with actual salt generation logic
		HashedKey:  "some-hashed-key",            // Replace with actual key hashing logic
		LastLogin:  time.Now(),
		VaultId:    1, // Replace with actual vault ID logic if needed
	}

	_, err := db.Exec("INSERT INTO users (username, client_salt, server_salt, hashed_key, last_login, vault_id) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username, user.ClientSalt, user.ServerSalt, user.HashedKey, user.LastLogin, user.VaultId)

	return err
}

func
