package absorber

import (
	"errors"
	"testing"

	"regenbrake/internal/model"
)

// engagingDriver 仅在前 N 次调用 Engage 时返回错误，之后成功。
// 用于模拟"投入失败→重试→成功"的场景。
type engagingDriver struct {
	failCount int
	calls     int
}

func (e *engagingDriver) Engage() error {
	e.calls++
	if e.calls <= e.failCount {
		return errors.New("engage failed")
	}
	return nil
}

func (e *engagingDriver) Discharge() error { return nil }

func TestEngageFailureReportsAndStaysIdle(t *testing.T) {
	// 用一个永远失败的驱动，验证投入失败时：
	//   1) 错误如实返回（不被压下）；
	//   2) 设备状态不得进入 absorbing（否则系统会显示"已投入"而装置并不吸收）；
	//   3) 故障原因与计数被记录，可供上报。
	d := NewDevice("A01", FaultyDriver{})

	err := d.Engage()
	if err == nil {
		t.Fatal("expected engage failure to be reported, got nil")
	}
	if d.State() == model.AbsorberAbsorbing {
		t.Fatalf("device must not show absorbing after failed engage, got %s", d.State())
	}
	if d.Status() == model.StatusOnline && d.State() != model.AbsorberIdle {
		t.Fatalf("device state should remain idle to allow retry, got status=%s state=%s", d.Status(), d.State())
	}
	if d.FaultCount() == 0 {
		t.Fatal("fault count must be recorded on engage failure")
	}
	if d.FaultReason() == "" {
		t.Fatal("fault reason must be recorded on engage failure")
	}
	if d.Metrics().Faults == 0 {
		t.Fatal("fault metric must be incremented on engage failure")
	}
}

func TestEngageRetryThenSucceeds(t *testing.T) {
	// 前 2 次投入失败、第 3 次成功：失败均如实返回，最终状态为已投入。
	drv := &engagingDriver{failCount: 2}
	d := NewDevice("A01", drv)

	var first, second error
	if err := d.Engage(); err != nil {
		first = err
	}
	if err := d.Engage(); err != nil {
		second = err
	}
	if first == nil || second == nil {
		t.Fatal("first two engage attempts must report failure")
	}
	if d.State() == model.AbsorberAbsorbing {
		t.Fatal("device must not show absorbing during failed attempts")
	}

	// 第三次成功。
	if err := d.Engage(); err != nil {
		t.Fatalf("third engage should succeed, got %v", err)
	}
	if d.State() != model.AbsorberAbsorbing {
		t.Fatalf("after successful engage, state should be absorbing, got %s", d.State())
	}
	if d.Metrics().Engaged == 0 {
		t.Fatal("engage metric must be incremented on success")
	}
}

func TestEngageSuccessClearsFaultArtifacts(t *testing.T) {
	// 先失败一次（留下故障记录），再成功投入。
	// 成功后进入 absorbing，故障历史仍保留供追溯（不抹除），但当前状态正确。
	drv := &engagingDriver{failCount: 1}
	d := NewDevice("A01", drv)

	if err := d.Engage(); err == nil {
		t.Fatal("first engage should fail")
	}
	if err := d.Engage(); err != nil {
		t.Fatalf("second engage should succeed, got %v", err)
	}
	if d.State() != model.AbsorberAbsorbing {
		t.Fatalf("state should be absorbing after success, got %s", d.State())
	}
}
