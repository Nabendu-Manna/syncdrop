package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ServerHandler struct {
	StorageDir string
}

type UploadResponse struct {
	Filename string `json:"filename"`
	Bytes    int64  `json:"bytes"`
	SHA256   string `json:"sha256"`
	Status   string `json:"status"`
}

func (h *ServerHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read multipart form stream
	mr, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "Failed to read multipart stream: "+err.Error(), http.StatusBadRequest)
		return
	}

	var uploadedFiles []UploadResponse

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "Error reading part: "+err.Error(), http.StatusBadRequest)
			return
		}

		filename := part.FileName()
		if filename == "" {
			continue
		}

		// Security: Prevent directory traversal
		cleanFilename := filepath.Base(filepath.Clean(filename))
		targetPath := filepath.Join(h.StorageDir, cleanFilename)

		// Verify path stays strictly within storage directory
		absTarget, _ := filepath.Abs(targetPath)
		absStorage, _ := filepath.Abs(h.StorageDir)
		if !strings.HasPrefix(absTarget, absStorage) {
			http.Error(w, "Invalid target path", http.StatusForbidden)
			return
		}

		// Temporary part file to prevent partial uploads corrupting storage
		tempPath := targetPath + ".part"
		dst, err := os.Create(tempPath)
		if err != nil {
			http.Error(w, "Failed to create target file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		hasher := sha256.New()
		multiWriter := io.MultiWriter(dst, hasher)

		written, err := io.Copy(multiWriter, part)
		dst.Close()

		if err != nil {
			os.Remove(tempPath)
			http.Error(w, "Write failure: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Atomically rename .part file to target filename
		if err := os.Rename(tempPath, targetPath); err != nil {
			os.Remove(tempPath)
			http.Error(w, "Failed to finalize file: "+err.Error(), http.StatusInternalServerError)
			return
		}

		uploadedFiles = append(uploadedFiles, UploadResponse{
			Filename: cleanFilename,
			Bytes:    written,
			SHA256:   hex.EncodeToString(hasher.Sum(nil)),
			Status:   "success",
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(uploadedFiles)
}

func (h *ServerHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"online"}`))
}
