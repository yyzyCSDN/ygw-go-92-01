package model

import "errors"

var (
	ErrTrainNotFound    = errors.New("train not found")
	ErrDeviceFault      = errors.New("device is in fault state")
	ErrAbsorbTimeout    = errors.New("absorb confirm timeout")
	ErrCapacityExceeded = errors.New("capacity exceeded")
	ErrRecordClosed     = errors.New("record journal closed")
	ErrSnapshotStale    = errors.New("snapshot is stale")
)
