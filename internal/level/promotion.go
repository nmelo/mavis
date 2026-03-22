package level

const AccuracyThreshold = 95.0

func PhasePassed(accuracy float64) bool {
	return accuracy >= AccuracyThreshold
}

func NextPhase(currentPhase string, lvl Level) string {
	switch currentPhase {
	case "chars":
		if lvl.HasWordDrills {
			return "words"
		}
		if lvl.HasCodeDrills {
			return "code"
		}
		return ""
	case "words":
		if lvl.HasCodeDrills {
			return "code"
		}
		return ""
	case "code":
		return ""
	}
	return ""
}

func PhasesFor(lvl Level) []string {
	phases := []string{"chars"}
	if lvl.HasWordDrills {
		phases = append(phases, "words")
	}
	if lvl.HasCodeDrills {
		phases = append(phases, "code")
	}
	return phases
}
