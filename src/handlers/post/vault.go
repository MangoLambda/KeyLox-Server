package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MangoLambda/KeyLox-Server/src/models"
)

// TODO: Fix docs
// @Summary Upload a vault
// @Description Upload a vault
// @Tags vault
// @Accept json
// @Produce json
// @Param user body models.RegisterRequest true "User"
// @Success 200 {object} string
// @Failure 409 {object} string
// @Failure 500 {object} string
// @Router /register [post]
func VaultHandler(db *sql.DB) http.HandlerFunc {
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
