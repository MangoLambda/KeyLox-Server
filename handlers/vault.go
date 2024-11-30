package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/MangoLambda/KeyLox-Server/models"
	"github.com/go-chi/chi/v5"
)

// @Summary Get a vault by ID
// @Description Get a vault by ID
// @Tags vaults
// @Produce json
// @Param username path string true "username"
// @Success 200 {object} models.VaultResponse
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Failure 500 {object} string
// @Router /vault/{id} [get]
func GetVault(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		usernameStr := chi.URLParam(r, "username")

		// TODO: Sanitze username
		// if err != nil {
		// 	http.Error(w, "Invalid Username", http.StatusBadRequest)
		// 	return
		// }

		var vault models.DBVault
		var vaultID int
		var createdAt string
		err := db.QueryRow("SELECT vault_id FROM users WHERE username = ?", usernameStr).Scan(&vaultID)
		if err == sql.ErrNoRows {
			http.Error(w, "Vault not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db.QueryRow("SELECT created_at FROM vaults WHERE id = ?", vaultID).Scan(&createdAt)
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
