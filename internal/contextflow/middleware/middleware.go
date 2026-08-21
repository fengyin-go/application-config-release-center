package middleware

import (
	"context"
	"time"
)

// WithDeadline 基于 parent 派生一个带超时的 context 并返回派生后的 context，
// 而非 context.Background()，否则调用方拿到的 context 永远不会到期，
// 超时状态会在这一层丢失。cancel 用于提前取消，调用方应通过 defer 调用。
func WithDeadline(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
