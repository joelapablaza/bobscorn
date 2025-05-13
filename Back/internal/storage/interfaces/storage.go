package interfaces

import "context"

type RateLimitStorage interface {
	CheckAndRecordRequest(ctx context.Context, clientIP string) (allowed bool, err error)
	CleanupOldRequests(ctx context.Context) error
}
