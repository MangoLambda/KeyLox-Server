package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MangoLambda/KeyLox-Server/auth"
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
func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerRequest models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validateInputs(registerRequest); err != nil {
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

	b64HashedKey, b64ServerSalt, err := auth.HashNewKey(registerRequest.Key)
	if err != nil {
		return err
	}

	// Set server salt, vault id, and last login
	user := models.DBUser{
		Username:   registerRequest.Username,
		ClientSalt: registerRequest.ClientSalt,
		ServerSalt: b64ServerSalt,
		HashedKey:  b64HashedKey,
		LastLogin:  time.Now(),
	}

	_, err = db.Exec("INSERT INTO users (username, client_salt, server_salt, hashed_key, last_login) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.ClientSalt, user.ServerSalt, user.HashedKey, user.LastLogin)

	return err
}

func validateInputs(registerRequest models.RegisterRequest) error {
	if registerRequest.Username == "" {
		return &models.InvalidInputError{Message: "Username cannot be empty"}
	}
	if registerRequest.Key == "" {
		return &models.InvalidInputError{Message: "Key cannot be empty"}
	}
	if registerRequest.ClientSalt == "" {
		return &models.InvalidInputError{Message: "Client salt cannot be empty"}
	}
	if _, err := base64.StdEncoding.DecodeString(registerRequest.Key); err != nil {
		return &models.InvalidInputError{Message: "Key must be a valid base64 string"}
	}
	if _, err := base64.StdEncoding.DecodeString(registerRequest.ClientSalt); err != nil {
		return &models.InvalidInputError{Message: "Client salt must be a valid base64 string"}
	}
	return nil
}
