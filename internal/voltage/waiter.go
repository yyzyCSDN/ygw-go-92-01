package voltage

import (
	"context"
	"errors"
	"time"

	"regenbrake/internal/model"
)

// WaitAbsorbConfirm 等待吸收装置的投入确认信号，最多等待 timeout。
//
// ctx 取消或等待超时都会中止等待并复位电压控制（脱离 absorbing、恢复逆变回馈），
// 使控制可在下一采样周期重新评估与动作。等待超时返回 model.ErrAbsorbTimeout，
// 外部取消（CancelWait 或 ctx 取消）返回对应的 ctx.Err()。
func (c *Controller) WaitAbsorbConfirm(ctx context.Context, timeout time.Duration) error {
	c.mu.Lock()
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	c.cancelWait = cancel
	c.mu.Unlock()
	defer cancel()
	defer c.clearWaitRegistration()

	select {
	case <-c.abs.ConfirmSignal():
		return nil
	case <-waitCtx.Done():
		// 等待超时或被外部取消：复位控制以脱离 absorbing 阻塞，
		// 避免控制卡在等待中导致电压响应迟缓，允许下一周期重新动作。
		c.SwitchToRestoring()
		if errors.Is(waitCtx.Err(), context.DeadlineExceeded) {
			return model.ErrAbsorbTimeout
		}
		return waitCtx.Err()
	}
}

// clearWaitRegistration 注销本次等待注册的取消函数。它持锁复用 cancelPendingWait，
// 避免外部 CancelWait 在本次等待结束后误取消下一次等待。
func (c *Controller) clearWaitRegistration() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cancelPendingWait()
}
