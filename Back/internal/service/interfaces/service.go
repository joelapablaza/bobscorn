package interfaces

import "context"

type CornService interface {
	CanBuyCorn(ctx context.Context, clientIP string) (allowed bool, err error)
}
