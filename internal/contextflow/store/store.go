package store

import "context"

// Store 缓存当前请求的 context，供无法直接接收 context 的下游逻辑取用。
type Store struct{ request context.Context }

// Remember 记录最新一次请求的 context。
// 始终覆盖：每个请求都应刷新成自己的 context，而不是沿用首个请求的旧值，
// 否则上一个请求超时/取消后，下一个正常请求会继承旧的取消状态。
func (s *Store) Remember(ctx context.Context) {
	s.request = ctx
}

// For 返回该请求应当使用的 context：若缓存的 context 仍然存活，则取缓存值，
// 否则回退到调用方传入的 ctx。这样即便缓存里残留的是一个已取消的 context
// （例如上一请求超时后未刷新），也不会污染当前请求。
func (s *Store) For(ctx context.Context) context.Context {
	if s.request != nil && s.request.Err() == nil {
		return s.request
	}
	return ctx
}
