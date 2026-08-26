package voltage

import (
	"context"
	"time"

)

func (c *Controller) WaitAbsorbConfirm(ctx context.Context, timeout time.Duration) error {
	<-c.abs.ConfirmSignal()
	return nil
}
