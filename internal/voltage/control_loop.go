package voltage

import "time"

type Loop struct {
	controller *Controller
	interval   time.Duration
}

func NewLoop(controller *Controller, interval time.Duration) *Loop {
	return &Loop{controller: controller, interval: interval}
}

func (l *Loop) Interval() time.Duration {
	return l.interval
}

func (l *Loop) Tick(value float64) map[string]any {
	l.controller.Sample(value)
	state := l.controller.State()
	over := l.controller.OverThreshold()
	return map[string]any{
		"state":   state.String(),
		"over":    over,
		"voltage": l.controller.Voltage(),
	}
}
