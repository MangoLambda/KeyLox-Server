package handlers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/sha3"
)

const maxUploadSize = 2 << 20  // 2 MB
const allocationSize = 3 << 20 // 3 MB

// @Summary Upload a vault
// @Description Upload a vault
// @Tags vault
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param username path string true "Username"
// @Success 200 {object} string
// @Failure 400 {object} string "Unable to parse form or file too big"
// @Failure 401 {object} string "Unauthorized"
// @Failure 409 {object} string
// @Failure 500 {object} string
// @Router /vault/{username} [post]
func VaultHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Basic Authentication
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
		err := r.ParseMultipartForm(allocationSize)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from form data
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Check the file size
		if handler.Size > maxUploadSize {
			http.Error(w, "The file is too big. 2 MB max.", http.StatusBadRequest)
			return
		}

		// TODO: put in a different package
		// Compute the SHA3-512 hash of the file
		hasher := sha3.New512()
		if _, err := io.Copy(hasher, file); err != nil {
			http.Error(w, "Unable to compute hash", http.StatusInternalServerError)
			return
		}
		hashSum := hasher.Sum(nil)
		hashFilename := fmt.Sprintf("%x", hashSum)

		// Reset the file pointer to the beginning
		file.Seek(0, io.SeekStart)

		// Define the directory path within the user's home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, "Unable to determine home directory", http.StatusInternalServerError)
			return
		}
		vaultDir := filepath.Join(homeDir, "keylox/vaults")

		// Create the directory if it doesn't exist
		if _, err := os.Stat(vaultDir); os.IsNotExist(err) {
			err = os.MkdirAll(vaultDir, os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create directory", http.StatusInternalServerError)
				return
			}
		}

		// Create a file on disk
		dst, err := os.Create(filepath.Join(vaultDir, hashFilename))
		if err != nil {
			http.Error(w, "Unable to create the file on disk", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the created file on disk
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
	}
}

func authorizeUser(db *sql.DB, username string, auth string) bool {
	// Check if the Authorization header is in the correct format
	const basicPrefix = "Basic "
	if !strings.HasPrefix(auth, basicPrefix) {
		return false
	}

	// Decode the base64-encoded credentials
	encodedCredentials := strings.TrimPrefix(auth, basicPrefix)
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return false
	}

	// Split the credentials into username and password
	credentials := strings.SplitN(string(decodedCredentials), ":", 2)
	if len(credentials) != 2 {
		return false
	}

	authUsername, authPassword := credentials[0], credentials[1]
	if authUsername != username {
		return false
	}

	if authUsername != "your-username" {

		return false
	}

	if authPassword != "your-password" {
		return false
	}

	return true
}
