package service

import (
	"context"
	"fmt"

	serviceinterfaces "bobscorn/internal/service/interfaces"
	storageinterfaces "bobscorn/internal/storage/interfaces"
)

type cornServiceImpl struct {
	storage storageinterfaces.RateLimitStorage
}

func NewCornService(storage storageinterfaces.RateLimitStorage) serviceinterfaces.CornService {
	return &cornServiceImpl{
		storage: storage,
	}
}

func (s *cornServiceImpl) CanBuyCorn(ctx context.Context, clientIP string) (bool, error) {
	allowed, err := s.storage.CheckAndRecordRequest(ctx, clientIP)
	if err != nil {
		return false, fmt.Errorf("error checking/recording request in storage: %w", err)
	}
	return allowed, nil
}
