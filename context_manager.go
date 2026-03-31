package logging

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ContextManager struct {
	defaultTimeout time.Duration
}

func NewContextManager(defaultTimeout time.Duration) *ContextManager {
	return &ContextManager{defaultTimeout: defaultTimeout}
}

func (cm *ContextManager) WithTimeout(parent *fiber.Ctx, timeout time.Duration, serviceName string) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = cm.defaultTimeout
	}

	baseCtx := parent.UserContext()
	if baseCtx == nil {
		baseCtx = parent.Context()
	}

	ctx, cancel := context.WithTimeout(baseCtx, timeout)
	return SetServiceName(ctx, serviceName), cancel
}

func (cm *ContextManager) WithDefaultTimeout(parent *fiber.Ctx, serviceName string) (context.Context, context.CancelFunc) {
	baseCtx := parent.UserContext()
	if baseCtx == nil {
		baseCtx = parent.Context()
	}

	ctx, cancel := context.WithTimeout(baseCtx, cm.defaultTimeout)
	return SetServiceName(ctx, serviceName), cancel
}

func (cm *ContextManager) WithDefaultTimeoutBasic() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), cm.defaultTimeout)
}

var DefaultContextManager = NewContextManager(10 * time.Second)

func WithTimeout(parent *fiber.Ctx, timeout time.Duration, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithTimeout(parent, timeout, serviceName)
}

func WithDefaultTimeout(parent *fiber.Ctx, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithDefaultTimeout(parent, serviceName)
}

func BackgroundWithTimeout(parent *fiber.Ctx, timeout time.Duration, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithTimeout(parent, timeout, serviceName)
}

func BackgroundWithDefaultTimeout(parent *fiber.Ctx, serviceName string) (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithDefaultTimeout(parent, serviceName)
}

func BackgroundWithDefaultTimeoutBasic() (context.Context, context.CancelFunc) {
	return DefaultContextManager.WithDefaultTimeoutBasic()
}
