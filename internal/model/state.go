package model

type AbsorberState string

const (
	AbsorberIdle        AbsorberState = "idle"
	AbsorberAbsorbing   AbsorberState = "absorbing"
	AbsorberDischarging AbsorberState = "discharging"
	AbsorberFault       AbsorberState = "fault"
)

type VoltageState string

const (
	VoltageNormal    VoltageState = "normal"
	VoltageHigh      VoltageState = "high"
	VoltageAbsorbing VoltageState = "absorbing"
	VoltageRestoring VoltageState = "restoring"
)

type CoopState string

const (
	CoopPlanning   CoopState = "planning"
	CoopAllocating CoopState = "allocating"
	CoopSettled    CoopState = "settled"
)

type DeviceStatus string

const (
	StatusOnline    DeviceStatus = "online"
	StatusFault     DeviceStatus = "fault"
	StatusRecovered DeviceStatus = "recovered"
)

func (s AbsorberState) String() string { return string(s) }

func (s VoltageState) String() string { return string(s) }

func (s CoopState) String() string { return string(s) }

func (s DeviceStatus) String() string { return string(s) }
