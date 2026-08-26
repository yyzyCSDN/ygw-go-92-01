package train

import "regenbrake/internal/model"

func (r *Registry) Braking(trainID string) (bool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	current, ok := r.trains[trainID]
	if !ok {
		return false, false
	}
	return current.Braking, true
}

func (r *Registry) BrakingTrains() []model.TrainState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]model.TrainState, 0)
	for _, current := range r.trains {
		if current.Braking {
			out = append(out, *current)
		}
	}
	return out
}
