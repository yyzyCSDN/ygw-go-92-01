package alarm

type Level string

const (
	LevelInfo     Level = "info"
	LevelWarning  Level = "warning"
	LevelCritical Level = "critical"
)

type Escalation struct {
	levels map[string]Level
}

func NewEscalation() *Escalation {
	return &Escalation{levels: make(map[string]Level)}
}

func (e *Escalation) Set(deviceID string, level Level) {
	e.levels[deviceID] = level
}

func (e *Escalation) Get(deviceID string) Level {
	if level, ok := e.levels[deviceID]; ok {
		return level
	}
	return LevelInfo
}

func (e *Escalation) Critical() []string {
	out := make([]string, 0)
	for id, level := range e.levels {
		if level == LevelCritical {
			out = append(out, id)
		}
	}
	return out
}
