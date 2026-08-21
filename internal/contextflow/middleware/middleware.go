package middleware

import (
	"context"
	"time"
)

func WithDeadline(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	_, cancel := context.WithTimeout(parent, timeout)
	return context.Background(), cancel
}
