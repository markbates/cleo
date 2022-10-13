package cleo

import (
	"context"
	"sync"
	"time"
)

// Context is a concurrent safe context.Context
// implementation that allows you to set an error
// on the context.
type Context struct {
	context.Context

	error  error
	cancel context.CancelFunc

	mu   sync.RWMutex
	once sync.Once
}

// Err returns the error set on the context.
func (c *Context) Err() error {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.error != nil {
		return c.error
	}

	if c.Context == nil {
		return nil
	}

	return c.Context.Err()
}

// SetError sets the error on the context.
func (c *Context) SetErr(err error) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.error = err
}

// Cancel cancels the context.
func (c *Context) Cancel() {
	if c == nil {
		return
	}

	c.once.Do(func() {
		c.mu.Lock()
		defer c.mu.Unlock()

		if c.cancel != nil {
			c.cancel()
		}
	})
}

// NewContext returns a new Context wrapping the given context.
// If the given context is canceled, the returned context will
// also be canceled.
func NewContext(ctx context.Context) (*Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-ctx.Done()
		cancel()
	}()

	cltx := &Context{
		Context: ctx,
		cancel:  cancel,
	}

	return cltx, cancel
}

// NewContextWithTimeout returns a new Context wrapping the given context with a timeout.
func ContextWithTimeout(ctx context.Context, timeout time.Duration) (*Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)

	go func() {
		<-ctx.Done()
		cancel()
	}()

	return &Context{
		Context: ctx,
		cancel:  cancel,
	}, cancel
}
