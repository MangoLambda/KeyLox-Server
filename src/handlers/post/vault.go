package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/sha3"
)

const maxUploadSize = 2 << 20  // 2 MB
const allocationSize = 3 << 20 // 3 MB

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

		if handler.Size > maxUploadSize {
			http.Error(w, "The file is too big. 2 MB max.", http.StatusBadRequest)
			return
		}

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

		// Create a file on disk
		dst, err := os.Create("/vaults/" + hashFilename)
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
