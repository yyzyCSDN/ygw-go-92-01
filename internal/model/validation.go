package model

func ValidTrainID(id string) bool {
	if len(id) == 0 || len(id) > 32 {
		return false
	}
	for _, ch := range id {
		if !(ch >= 'a' && ch <= 'z') && !(ch >= 'A' && ch <= 'Z') && !(ch >= '0' && ch <= '9') && ch != '-' && ch != '_' {
			return false
		}
	}
	return true
}

func ValidCapacity(capacity float64) bool {
	return capacity >= 0
}

func ValidVoltage(voltage float64) bool {
	return voltage >= 0 && voltage <= 3000
}
