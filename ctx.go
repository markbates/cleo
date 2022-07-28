package cleo

import (
	"context"
	"sync"
	"time"
)

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

type Context struct {
	context.Context

	Error  error
	cancel context.CancelFunc

	mu   sync.RWMutex
	once sync.Once
}

func (c *Context) Err() error {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.Error != nil {
		return c.Error
	}

	if c.Context == nil {
		return nil
	}

	return c.Context.Err()
}

func (c *Context) SetErr(err error) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.Error = err
}

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
