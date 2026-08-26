package inverter

import "sync"

type Inverter struct {
	mu       sync.Mutex
	active   bool
	feedback float64
	ramp     *Ramp
}

func New() *Inverter {
	return &Inverter{ramp: NewRamp(10)}
}

func (i *Inverter) Enable() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.active = true
}

func (i *Inverter) Disable() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.active = false
}

func (i *Inverter) Active() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.active
}

func (i *Inverter) SetFeedback(power float64) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.feedback = power
	i.ramp.SetTarget(power)
}

func (i *Inverter) FeedbackPower() float64 {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.feedback
}

func (i *Inverter) Step() float64 {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.feedback = i.ramp.Step()
	return i.feedback
}

func (i *Inverter) Settled() bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.ramp.Settled()
}

func (i *Inverter) Current() float64 {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.ramp.Current()
}
