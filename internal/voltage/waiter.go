package voltage

import (
	"context"
	"time"

	"regenbrake/internal/model"
)

func (c *Controller) WaitAbsorbConfirm(ctx context.Context, timeout time.Duration) error {
	child, cancel := context.WithTimeout(ctx, timeout)
	c.mu.Lock()
	c.cancelWait = cancel
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		if c.cancelWait != nil {
			c.cancelWait()
			c.cancelWait = nil
		}
		c.mu.Unlock()
	}()
	select {
	case <-c.abs.ConfirmSignal():
		return nil
	case <-child.Done():
		return model.ErrAbsorbTimeout
	}
}
