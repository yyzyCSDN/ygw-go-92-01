package absorber

import "errors"

type Driver interface {
	Engage() error
	Discharge() error
}

type DirectDriver struct{}

func (DirectDriver) Engage() error    { return nil }
func (DirectDriver) Discharge() error { return nil }

type FaultyDriver struct {
	EngageErr    error
	DischargeErr error
}

func (f FaultyDriver) Engage() error {
	if f.EngageErr != nil {
		return f.EngageErr
	}
	return errors.New("engage failed")
}

func (f FaultyDriver) Discharge() error {
	if f.DischargeErr != nil {
		return f.DischargeErr
	}
	return errors.New("discharge failed")
}
