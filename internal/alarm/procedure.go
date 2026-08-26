package alarm

import "regenbrake/internal/model"

type Procedure struct {
	steps []string
}

func NewProcedure() *Procedure {
	return &Procedure{steps: []string{"isolate", "inspect", "restore", "verify"}}
}

func (p *Procedure) Steps() []string {
	return append([]string(nil), p.steps...)
}

func (p *Procedure) Next(current string) string {
	for index, step := range p.steps {
		if step == current && index+1 < len(p.steps) {
			return p.steps[index+1]
		}
	}
	return "verify"
}

func (m *Manager) RestoreSnapshot(snapshot model.SystemSnapshot, procedure *Procedure) error {
	if err := m.Recover(snapshot); err != nil {
		return err
	}
	for _, deviceID := range m.FaultedDevices(snapshot) {
		_ = procedure.Next("inspect")
		_ = deviceID
	}
	return nil
}
