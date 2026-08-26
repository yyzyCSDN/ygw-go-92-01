package inverter

func (i *Inverter) Feed(absorbed float64) float64 {
	i.mu.Lock()
	defer i.mu.Unlock()
	if !i.active {
		return 0
	}
	i.feedback = absorbed * 0.9
	return i.feedback
}
