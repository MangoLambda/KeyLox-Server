package handlers

import (
	"database/sql"
	"encoding/base32"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	authorization "github.com/MangoLambda/KeyLox-Server/src/auth"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/sha3"
)

const maxUploadSize = 2 << 20  // 2 MB
const allocationSize = 3 << 20 // 3 MB

// @Summary Upload a vault
// @Description Upload a vault
// @Tags vault
// @Accept multipart/form-data
// @Produce json
// @Param username path string true "Username"
// @Param file formData file true "File to upload"
// @Success 200 {object} string
// @Failure 400 {object} string "Unable to parse form or file too big"
// @Failure 401 {object} string "Unauthorized"
// @Failure 409 {object} string
// @Failure 500 {object} string
// @Router /vault/{username} [post]
func VaultHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isAuthorized, err := isUserAuthorized(db, r)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if !isAuthorized {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		file, err := getFileFromRequest(w, r)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Compute the SHA3-512 hash of the file
		var filename string
		filename, err = calculateFilename(file)
		if err != nil {
			http.Error(w, "Unable to compute filename", http.StatusInternalServerError)
		}

		var dst *os.File
		dst, err = createFileOnDisk(filename)
		if err != nil {
			http.Error(w, "Unable to create the file on disk", http.StatusInternalServerError)
		}
		defer dst.Close()

		// Copy the uploaded file to the created file on disk
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Unable to save the file", http.StatusInternalServerError)
			return
		}

		// TODO Update the vault table in the DB

		// Respond to the client
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File uploaded successfully"))
	}
}

func isUserAuthorized(db *sql.DB, r *http.Request) (bool, error) {
	routeUsername := chi.URLParam(r, "username")
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false, fmt.Errorf("no authorization header provided")
	}

	// Retrieve the hashed password and server salt from the database
	var storedPasswordHash, storedServerSalt string
	err := db.QueryRow("SELECT hashed_key, server_salt FROM users WHERE username = ?", routeUsername).Scan(&storedPasswordHash, &storedServerSalt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("username not found: %v", routeUsername)
		}
		return false, fmt.Errorf("database query error: %v", err)
	}

	areCredentialsValid, authHeaderUsername, err := authorization.VerifyCredentials(storedPasswordHash, storedServerSalt, auth)
	if err != nil {
		return false, fmt.Errorf("unable to verify credentials: %v", err)
	}

	// Check if the username in the route matches the username in the credentials
	if authHeaderUsername != routeUsername {
		return false, fmt.Errorf("username in the route does not match the username in the credentials [%v, %v]",
			routeUsername, authHeaderUsername)
	}

	return areCredentialsValid, nil
}

func getFileFromRequest(w http.ResponseWriter, r *http.Request) (file multipart.File, err error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	err = r.ParseMultipartForm(allocationSize)
	if err != nil {
		return nil, fmt.Errorf("unable to parse form: %v", err)
	}

	// Retrieve the file from form data
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("error retrieving the file: %v", err)
	}

	// Check the file size
	if handler.Size > maxUploadSize {
		file.Close()
		return nil, fmt.Errorf("the file is too big. 2 MB max: %v", err)
	}

	return file, nil
}

func createFileOnDisk(filename string) (dst *os.File, err error) {
	// Define the directory path within the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to create directory: %v", err)
	}
	vaultDir := filepath.Join(homeDir, "keylox/vaults")

	// Create the directory if it doesn't exist
	if _, err = os.Stat(vaultDir); os.IsNotExist(err) {
		err = os.MkdirAll(vaultDir, os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("unable to create directory: %v", err)

		}
	}

	// Create a file on disk
	dst, err = os.Create(filepath.Join(vaultDir, filename))
	if err != nil {
		return nil, fmt.Errorf("unable to create the file on disk: %v", err)
	}

	return dst, nil
}

func calculateFilename(file multipart.File) (string, error) {
	// Compute the SHA3-512 hash of the file
	hasher := sha3.New512()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("unable to compute hash: %v", err)
	}
	hashSum := hasher.Sum(nil)
	base32Filename := base32.StdEncoding.EncodeToString(hashSum)

	// Reset the file pointer to the beginning
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("unable to reset the file pointer: %v", err)
	}

	return base32Filename, nil
}
