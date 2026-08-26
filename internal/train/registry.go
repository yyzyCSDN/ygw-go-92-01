package train

func (r *Registry) Exists(trainID string) bool {
	_, ok := r.Lookup(trainID)
	return ok
}

func (r *Registry) Remove(trainID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.trains, trainID)
}
