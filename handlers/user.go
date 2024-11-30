package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MangoLambda/KeyLox-Server/models"
	"github.com/go-chi/chi/v5"
)

// @Summary Gets a specific user
// @Description Gets a specific user
// @Tags user
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} models.UserResponse
// @Success 404 {object} string
// @Failure 500 {object} string
// @Router /user/{username} [get]
func GetUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := chi.URLParam(r, "username")

		var user models.DBUser
		rows, err := db.Query("SELECT id, username, client_salt FROM users WHERE username = ?", username)
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
			http.Error(w, "User not found", http.StatusNotFound)
		}
	}
}
