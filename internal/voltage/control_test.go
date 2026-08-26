package voltage

import (
	"errors"
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/inverter"
	"regenbrake/internal/model"
)

// repeatFailDriver 在前 fail 次调用 Engage 时返回错误，之后成功。
type repeatFailDriver struct {
	fail int
	n    int
}

func (r *repeatFailDriver) Engage() error {
	r.n++
	if r.n <= r.fail {
		return errors.New("engage failed")
	}
	return nil
}

func (r *repeatFailDriver) Discharge() error { return nil }

func newControllerWithDriver(drv absorber.Driver) (*Controller, *absorber.Device) {
	dev := absorber.NewDevice("A01", drv)
	inv := inverter.New()
	c := NewController(dev, inv, 900.0)
	return c, dev
}

// 模拟 SampleVoltage 的核心循环：电压超阈值则投入，失败时上报并重试，
// 全部失败则电压状态保持 high（不能显示已投入）。
func engageWithRetry(c *Controller, maxAttempts int) error {
	var last error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := c.SwitchToAbsorbing(); err == nil {
			return nil
		} else {
			last = err
		}
	}
	return last
}

func TestSwitchToAbsorbingFailureKeepsStateHigh(t *testing.T) {
	// 驱动永远失败：投入失败后电压状态必须保持 high，不能进入 absorbing。
	c, dev := newControllerWithDriver(&absorber.FaultyDriver{})
	c.Sample(950.0) // 超阈值
	if !c.OverThreshold() {
		t.Fatal("expected voltage over threshold")
	}

	err := c.SwitchToAbsorbing()
	if err == nil {
		t.Fatal("expected engage failure to be reported, got nil")
	}
	if c.State() == model.VoltageAbsorbing {
		t.Fatal("voltage must not show absorbing when engage failed (system would lie 'engaged' while not absorbing)")
	}
	if c.State() != model.VoltageHigh {
		t.Fatalf("voltage state should remain high after failed engage, got %s", c.State())
	}
	if dev.State() == model.AbsorberAbsorbing {
		t.Fatal("device must not show absorbing after failed engage")
	}
	// 逆变器不应被关闭（投入未成功时 SwitchToAbsorbing 不得调用 Disable）。
	c.inv.Enable()
	if !c.inv.Active() {
		t.Fatal("inverter should remain active when engage failed")
	}
	if err := c.SwitchToAbsorbing(); err == nil {
		t.Fatal("second engage should also fail")
	}
	if !c.inv.Active() {
		t.Fatal("inverter must remain active across failed engages")
	}
}

func TestSwitchToAbsorbingRetryThenSucceeds(t *testing.T) {
	// 前 2 次失败、第 3 次成功：重试后状态正确进入 absorbing。
	c, dev := newControllerWithDriver(&repeatFailDriver{fail: 2})
	c.Sample(950.0)

	if err := engageWithRetry(c, 3); err != nil {
		t.Fatalf("retry should eventually succeed, got %v", err)
	}
	if c.State() != model.VoltageAbsorbing {
		t.Fatalf("after successful engage, voltage state should be absorbing, got %s", c.State())
	}
	if dev.State() != model.AbsorberAbsorbing {
		t.Fatalf("device should be absorbing after success, got %s", dev.State())
	}
}

func TestSwitchToAbsorbingAllRetriesFail(t *testing.T) {
	// 全部重试失败：返回最后一次错误，电压状态保持 high，不显示已投入。
	c, _ := newControllerWithDriver(&absorber.FaultyDriver{})
	c.Sample(950.0)

	err := engageWithRetry(c, 3)
	if err == nil {
		t.Fatal("expected error when all retries fail")
	}
	if c.State() == model.VoltageAbsorbing {
		t.Fatal("must not show absorbing when all retries failed")
	}
	if c.State() != model.VoltageHigh {
		t.Fatalf("voltage state should remain high, got %s", c.State())
	}
}
