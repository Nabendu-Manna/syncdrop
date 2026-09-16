package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/nabendu-manna/syncdrop/server/internal/api"
	"github.com/nabendu-manna/syncdrop/server/internal/discovery"
)

func main() {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080
	}

	storageDir := getEnv("STORAGE_DIR", "./storage")
	absStorageDir, err := filepath.Abs(storageDir)
	if err != nil {
		log.Fatalf("Failed to resolve storage directory: %v", err)
	}

	if err := os.MkdirAll(absStorageDir, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}

	serverID := getEnv("SERVER_ID", uuid.New().String())

	// 1. Start mDNS Advertiser
	mdnsServer, err := discovery.StartMDNS(port, serverID)
	if err != nil {
		log.Printf("[WARN] mDNS failed to initialize: %v", err)
	} else {
		defer mdnsServer.Shutdown()
	}

	// 2. Setup Routes
	handler := &api.ServerHandler{StorageDir: absStorageDir}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", handler.HealthHandler)
	mux.HandleFunc("/api/v1/upload", handler.UploadHandler)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      mux,
		ReadTimeout:  30 * time.Minute, // Allow large file streams
		WriteTimeout: 30 * time.Minute,
	}

	// 3. Graceful Shutdown listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[HTTP] SyncDrop server listening on http://0.0.0.0:%d", port)
		log.Printf("[Storage] Saving files to %s", absStorageDir)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited cleanly.")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
