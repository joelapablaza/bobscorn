package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	storageinterfaces "bobscorn/internal/storage/interfaces"
)

type PostgresRateLimitStorage struct {
	db              *sql.DB
	rateLimitWindow time.Duration
}

func NewPostgresRateLimitStorage(db *sql.DB) storageinterfaces.RateLimitStorage {
	rateLimitSecondsStr := os.Getenv("RATE_LIMIT_WINDOW_SECONDS")
	if rateLimitSecondsStr == "" {
		rateLimitSecondsStr = "60"
	}
	rateLimitSeconds, err := strconv.Atoi(rateLimitSecondsStr)
	if err != nil {
		log.Printf("Warning: Invalid RATE_LIMIT_WINDOW_SECONDS '%s', defaulting to 60. Error: %v", rateLimitSecondsStr, err)
		rateLimitSeconds = 60
	}
	return &PostgresRateLimitStorage{
		db:              db,
		rateLimitWindow: time.Duration(rateLimitSeconds) * time.Second,
	}
}

func (s *PostgresRateLimitStorage) CheckAndRecordRequest(ctx context.Context, clientIP string) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Ensure rollback if commit is not called
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			// err is non-nil; don't change it
			// if it was already rolled back, this is a no-op
			if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
				log.Printf("Transaction rollback failed: %v", rbErr)
			}
		} else {
			// err is nil; if an error occurred during a successful commit,
			// the transaction might still be open. This check is for safety.
			// However, if commit was successful, err would be nil, and this block isn't hit.
		}
	}()

	var lastRequestTime sql.NullTime

	// Lock the row for this client_ip to prevent race conditions from concurrent requests from the same IP.
	// If the row doesn't exist, QueryRowContext returns sql.ErrNoRows.
	query := "SELECT last_request_at FROM client_requests WHERE client_ip = $1 FOR UPDATE"
	err = tx.QueryRowContext(ctx, query, clientIP).Scan(&lastRequestTime)

	if err == sql.ErrNoRows {
		insertQuery := "INSERT INTO client_requests (client_ip, last_request_at) VALUES ($1, NOW())"
		_, err = tx.ExecContext(ctx, insertQuery, clientIP)
		if err != nil {
			return false, fmt.Errorf("failed to insert new request record: %w", err)
		}
		err = tx.Commit()
		if err != nil {
			return false, fmt.Errorf("failed to commit transaction for new request: %w", err)
		}
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to query last request time: %w", err)
	}

	if !lastRequestTime.Valid {
		// This case should ideally not happen if last_request_at is NOT NULL.
		// If there are some errors with the database, this can happen.
		// Treat as if it's a new request.
		updateQuery := "UPDATE client_requests SET last_request_at = NOW() WHERE client_ip = $1"
		_, err = tx.ExecContext(ctx, updateQuery, clientIP)
		if err != nil {
			return false, fmt.Errorf("failed to update request time for NULL last_request_at: %w", err)
		}
		err = tx.Commit()
		if err != nil {
			return false, fmt.Errorf("failed to commit transaction for NULL last_request_at update: %w", err)
		}
		return true, nil
	}

	// Check if the last request was within the rate limit window (e.g., 1 minute).
	if time.Since(lastRequestTime.Time) < s.rateLimitWindow {
		// Request is too soon. Rate limit exceeded.
		// We don't update the timestamp here. The transaction will be rolled back by defer.
		// Explicitly rolling back for clarity and to release the lock sooner.
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			log.Printf("Error on rollback for rate limited request: %v", rbErr)
		}
		return false, nil
	}

	updateQuery := "UPDATE client_requests SET last_request_at = NOW() WHERE client_ip = $1"
	_, err = tx.ExecContext(ctx, updateQuery, clientIP)
	if err != nil {
		return false, fmt.Errorf("failed to update last request time: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("failed to commit transaction for allowed update: %w", err)
	}
	return true, nil
}

// CleanupOldRequests removes records older than 5 minutes (or configured CLEANUP_INTERVAL_MINUTES).
// This helps keep the client_requests table from growing indefinitely.
func (s *PostgresRateLimitStorage) CleanupOldRequests(ctx context.Context) error {
	cleanupIntervalMinutesStr := os.Getenv("CLEANUP_INTERVAL_MINUTES")
	if cleanupIntervalMinutesStr == "" {
		cleanupIntervalMinutesStr = "5"
	}
	cleanupIntervalMinutes, convErr := strconv.Atoi(cleanupIntervalMinutesStr)
	if convErr != nil {
		log.Printf("Warning: Invalid CLEANUP_INTERVAL_MINUTES for cleanup logic '%s', defaulting to 5. Error: %v", cleanupIntervalMinutesStr, convErr)
		cleanupIntervalMinutes = 5
	}

	// Delete records where last_request_at is older than 'cleanupIntervalMinutes' ago.
	// This ensures we don't delete records that are still relevant for the current rate-limiting window.
	// For example, if rate limit is 1 min and cleanup is 5 mins, we delete records older than 5 mins.
	query := fmt.Sprintf("DELETE FROM client_requests WHERE last_request_at < NOW() - INTERVAL '%d minutes'", cleanupIntervalMinutes)
	result, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete old request records: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("CleanupOldRequests: Deleted %d old records.", rowsAffected)
	return nil
}
