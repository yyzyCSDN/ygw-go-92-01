package alarm

import "regenbrake/internal/model"

func (m *Manager) FaultedDevices(snapshot model.SystemSnapshot) []string {
	out := make([]string, 0)
	for _, snap := range snapshot.Devices {
		if snap.Status == model.StatusFault {
			out = append(out, snap.ID)
		}
	}
	return out
}
