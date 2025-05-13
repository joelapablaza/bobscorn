package utils

import (
	"context"
	"log"
	"time"

	storageinterfaces "bobscorn/internal/storage/interfaces"
)

func StartCleanupTask(ctx context.Context, storage storageinterfaces.RateLimitStorage, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Starting cleanup task for old requests every %v", interval)

	// Run once immediately at startup, then on ticker
	cleanup := func() {
		log.Println("Running cleanup task for old requests...")
		err := storage.CleanupOldRequests(context.Background())
		if err != nil {
			log.Printf("Error cleaning up old requests: %v", err)
		} else {
			log.Println("Cleanup task completed.")
		}
	}

	cleanup()

	for {
		select {
		case <-ticker.C:
			cleanup()
		case <-ctx.Done(): // Listen for cancellation signal
			log.Println("Stopping cleanup task...")
			return
		}
	}
}
