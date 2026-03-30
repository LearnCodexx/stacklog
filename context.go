package logging

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

// ContextManager stores default timeout configuration for request contexts.
type ContextManager struct {
	defaultTimeout time.Duration
}

// NewContextManager constructs a context manager with the provided default timeout.
func NewContextManager(defaultTimeout time.Duration) *ContextManager {
	return &ContextManager{
		defaultTimeout: defaultTimeout,
	}
}

// WithTimeout wraps the Fiber context with a timeout and service name metadata.
func (cm *ContextManager) WithTimeout(parent *fiber.Ctx, timeout time.Duration, serviceName string) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = cm.defaultTimeout
	}

	baseCtx := parent.UserContext()
	if baseCtx == nil {
		baseCtx = parent.Context()
	}

	ctx, cancel := context.WithTimeout(baseCtx, timeout)
	ctxWithService := SetServiceName(ctx, serviceName)

	return ctxWithService, cancel
}

// WithDefaultTimeout applies the default timeout using the provided Fiber context.
func (cm *ContextManager) WithDefaultTimeout(parent *fiber.Ctx, serviceName string) (context.Context, context.CancelFunc) {
	baseCtx := parent.UserContext()
	if baseCtx == nil {
		baseCtx = parent.Context()
	}

	ctx, cancel := context.WithTimeout(baseCtx, cm.defaultTimeout)
	ctxWithService := SetServiceName(ctx, serviceName)

	return ctxWithService, cancel
}

// WithDefaultTimeoutBasic creates a background context with the default timeout.
func (cm *ContextManager) WithDefaultTimeoutBasic() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), cm.defaultTimeout)
	return ctx, cancel
}

// DefaultContextManager is a shared instance with a 10s timeout.
var DefaultContextManager = NewContextManager(10 * time.Second)

// WithTimeout wraps the default manager helper for Fiber requests.
func WithTimeout(parent *fiber.Ctx, timeout time.Duration, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithTimeout(parent, timeout, serviceName)
}

// WithDefaultTimeout wraps the default manager for common requests.
func WithDefaultTimeout(parent *fiber.Ctx, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithDefaultTimeout(parent, serviceName)
}

// BackgroundWithTimeout is an alias to make intent clear for background Fiber contexts.
func BackgroundWithTimeout(parent *fiber.Ctx, timeout time.Duration, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithTimeout(parent, timeout, serviceName)
}

// BackgroundWithDefaultTimeout uses the default timeout for background Fibers.
func BackgroundWithDefaultTimeout(parent *fiber.Ctx, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithDefaultTimeout(parent, serviceName)
}

// BackgroundWithDefaultTimeoutBasic returns a background context with the default timeout.
func BackgroundWithDefaultTimeoutBasic() (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithDefaultTimeoutBasic()
}
