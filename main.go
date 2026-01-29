package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Image struct {
	ID           string    `json:"id"`
	Filename     string    `json:"filename"`
	OriginalName string    `json:"originalName"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mimeType"`
	UploadedAt   time.Time `json:"uploadedAt"`
	Path         string    `json:"path"`
}

var images []Image
var uploadsDir = "./uploads"
var dataDir = "./data"
var metadataFile = "./data/metadata.json"

func main() {
	// Create necessary directories
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Printf("Warning: Failed to create uploads directory: %v", err)
	}
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Warning: Failed to create data directory: %v", err)
	}

	// Load existing metadata
	loadMetadata()

	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/images", getAllImages).Methods("GET")
	router.HandleFunc("/api/images/{id}", getImage).Methods("GET")
	router.HandleFunc("/api/images/upload", uploadImage).Methods("POST")
	router.HandleFunc("/api/images/{id}", deleteImage).Methods("DELETE")
	router.HandleFunc("/api/images/file/{id}", serveImageFile).Methods("GET")

	// Serve static frontend files (SPA support)
	spa := spaHandler{staticPath: "./frontend/build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	fmt.Println("Server running on http://localhost:6000")
	log.Fatal(http.ListenAndServe(":6000", router))
}

// Load metadata from JSON file
func loadMetadata() {
	data, err := os.ReadFile(metadataFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, start with empty array
			images = []Image{}
			saveMetadata() // Create the file
			return
		}
		log.Printf("Error reading metadata: %v. Starting with empty array.", err)
		images = []Image{}
		return
	}

	if err := json.Unmarshal(data, &images); err != nil {
		log.Printf("Error parsing metadata: %v. Starting with empty array.", err)
		images = []Image{}
		return
	}
}

// Save metadata to JSON file
func saveMetadata() error {
	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		log.Printf("Error marshaling metadata: %v", err)
		return err
	}

	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		log.Printf("Error writing metadata: %v", err)
		return err
	}

	return nil
}

// Sanitize filename to keep only alphanumeric, dots, and hyphens
func sanitizeFilename(filename string) string {
	// Replace any character that's not alphanumeric, dot, or hyphen with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9.\-]`)
	return reg.ReplaceAllString(filename, "_")
}

// GET /api/images - Get all images sorted by upload date (newest first)
func getAllImages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Sort by UploadedAt descending (newest first)
	sortedImages := make([]Image, len(images))
	copy(sortedImages, images)
	sort.Slice(sortedImages, func(i, j int) bool {
		return sortedImages[i].UploadedAt.After(sortedImages[j].UploadedAt)
	})

	json.NewEncoder(w).Encode(sortedImages)
}

// GET /api/images/{id} - Get single image metadata
func getImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	for _, img := range images {
		if img.ID == id {
			json.NewEncoder(w).Encode(img)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Image not found"})
}

// POST /api/images/upload - Upload new image
func uploadImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse multipart form
	if err := r.ParseMultipartForm(0); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse form"})
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "No file provided"})
		return
	}
	defer file.Close()

	// Generate UUID and create filename
	imageID := uuid.New().String()
	sanitizedName := sanitizeFilename(header.Filename)
	filename := fmt.Sprintf("%s_%s", imageID, sanitizedName)
	filePath := filepath.Join(uploadsDir, filename)

	// Save file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save image"})
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read file"})
		return
	}

	// Detect MIME type
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Create Image object
	newImage := Image{
		ID:           imageID,
		Filename:     filename,
		OriginalName: header.Filename,
		Size:         size,
		MimeType:     mimeType,
		UploadedAt:   time.Now(),
		Path:         filePath,
	}

	// Add to array and persist
	images = append(images, newImage)
	if err := saveMetadata(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save metadata"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newImage)
}

// DELETE /api/images/{id} - Delete image
func deleteImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	// Find image
	found := false
	var imgToDelete Image
	for idx, img := range images {
		if img.ID == id {
			found = true
			imgToDelete = img
			// Remove from array
			images = append(images[:idx], images[idx+1:]...)
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Image not found"})
		return
	}

	// Delete physical file (ignore error if file doesn't exist)
	if err := os.Remove(imgToDelete.Path); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Failed to delete file %s: %v", imgToDelete.Path, err)
	}

	// Persist metadata
	if err := saveMetadata(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update metadata"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Image deleted successfully"})
}

// GET /api/images/file/{id} - Serve actual image file
func serveImageFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Find image metadata
	var img *Image
	for i := range images {
		if images[i].ID == id {
			img = &images[i]
			break
		}
	}

	if img == nil {
		http.NotFound(w, r)
		return
	}

	// Check if file exists
	if _, err := os.Stat(img.Path); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Set headers
	w.Header().Set("Content-Type", img.MimeType)
	w.Header().Set("Content-Disposition", "inline")

	// Serve file
	http.ServeFile(w, r, img.Path)
}

// SPA Handler for serving React frontend
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Build full path
	path := filepath.Join(h.staticPath, r.URL.Path)

	// Check if path exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// File doesn't exist, serve index.html for SPA routing
		indexFile := filepath.Join(h.staticPath, h.indexPath)
		if _, err := os.Stat(indexFile); os.IsNotExist(err) {
			// Build folder doesn't exist
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Frontend build not found. Run 'cd frontend && npm run build' first.")
			return
		}
		http.ServeFile(w, r, indexFile)
		return
	}

	// If path is a directory, serve index.html
	fi, err := os.Stat(path)
	if err == nil && fi.IsDir() {
		indexFile := filepath.Join(path, h.indexPath)
		if _, err := os.Stat(indexFile); err == nil {
			http.ServeFile(w, r, indexFile)
			return
		}
		// No index.html in directory, serve root index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	}

	// Serve the file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
