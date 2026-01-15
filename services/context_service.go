package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestContext struct {
	RequestID string
	UserID    uint
	TraceID   string
	StartTime time.Time
	Timeout   time.Duration
}

type AppContext struct {
	context.Context
	RequestID string
	UserID    uint
	TraceID   string
	StartTime time.Time
}

func NewRequestContext() *RequestContext {
	return &RequestContext{
		RequestID: uuid.New().String(),
		TraceID:   uuid.New().String(),
		StartTime: time.Now(),
		Timeout:   30 * time.Second,
	}
}

func WithRequestContext(ctx context.Context, reqCtx *RequestContext, userID uint) *AppContext {
	return &AppContext{
		Context:   ctx,
		RequestID: reqCtx.RequestID,
		UserID:    userID,
		TraceID:   reqCtx.TraceID,
		StartTime: reqCtx.StartTime,
	}
}

func (c *AppContext) WithTimeout(timeout time.Duration) *AppContext {
	ctx, cancel := context.WithTimeout(c.Context, timeout)

	// Store cancel function for potential cleanup
	go func() {
		<-ctx.Done()
		cancel()
	}()

	return &AppContext{
		Context:   ctx,
		RequestID: c.RequestID,
		UserID:    c.UserID,
		TraceID:   c.TraceID,
		StartTime: c.StartTime,
	}
}

func (c *AppContext) WithValue(key, value interface{}) *AppContext {
	ctx := context.WithValue(c.Context, key, value)
	return &AppContext{
		Context:   ctx,
		RequestID: c.RequestID,
		UserID:    c.UserID,
		TraceID:   c.TraceID,
		StartTime: c.StartTime,
	}
}

func (c *AppContext) GetElapsedTime() time.Duration {
	return time.Since(c.StartTime)
}

func (c *AppContext) IsExpired() bool {
	select {
	case <-c.Context.Done():
		return true
	default:
		return false
	}
}

// Context keys for type safety
type ContextKey string

const (
	ContextKeyRequestID ContextKey = "request_id"
	ContextKeyUserID    ContextKey = "user_id"
	ContextKeyTraceID   ContextKey = "trace_id"
)
