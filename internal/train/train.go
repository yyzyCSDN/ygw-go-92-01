package train

import (
	"sort"
	"sync"

	"regenbrake/internal/model"
)

type BrakingListener interface {
	OnTrainBrakingChange(trainID string, braking bool)
}

type Registry struct {
	mu        sync.RWMutex
	trains    map[string]*model.TrainState
	listeners []BrakingListener
	revisions map[string]int
}

func NewRegistry() *Registry {
	return &Registry{trains: make(map[string]*model.TrainState), revisions: make(map[string]int)}
}

func (r *Registry) Register(state model.TrainState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := state
	r.trains[state.ID] = &copy
}

func (r *Registry) SetBraking(trainID string, braking bool, power float64) error {
	r.mu.Lock()
	current, ok := r.trains[trainID]
	if !ok {
		r.mu.Unlock()
		return model.ErrTrainNotFound
	}
	current.Braking = braking
	current.Power = power
	r.revisions[trainID]++
	r.mu.Unlock()
	return nil
}

func (r *Registry) Lookup(trainID string) (*model.TrainState, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	current, ok := r.trains[trainID]
	if !ok {
		return nil, false
	}
	copy := *current
	return &copy, true
}

func (r *Registry) AddListener(listener BrakingListener) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners = append(r.listeners, listener)
}

func (r *Registry) List() []model.TrainState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.trains))
	for id := range r.trains {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]model.TrainState, 0, len(ids))
	for _, id := range ids {
		out = append(out, *r.trains[id])
	}
	return out
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.trains)
}

func (r *Registry) Revision(trainID string) int {
	return 0
}

func (r *Registry) BrakingOrFalse(trainID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	current, ok := r.trains[trainID]
	if !ok {
		return false
	}
	return current.Braking
}
