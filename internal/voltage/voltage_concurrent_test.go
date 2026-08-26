package voltage

import (
	"sync"
	"testing"

	"regenbrake/internal/absorber"
	"regenbrake/internal/inverter"
	"regenbrake/internal/model"
)

// reproduceConcurrentSwitch launches the two conflicting control instructions
// concurrently: "切吸收" (SwitchToAbsorbing) and "切回馈/恢复" (SwitchToRestoring /
// Restore).  Before the fix the absorber device state is written without the
// device mutex from SwitchTo while Engage holds it, and the controller fields
// (state/switches/inverter) are mutated without the controller mutex, producing
// torn state and -race reports.  The two paths must converge on a consistent
// final state regardless of interleaving.
func TestConcurrentAbsorbRestoreConsistency(t *testing.T) {
	abs := absorber.NewDevice("A01", absorber.DirectDriver{})
	inv := inverter.New()
	c := NewController(abs, inv, 900.0)

	var wg sync.WaitGroup
	const rounds = 200
	for r := 0; r < rounds; r++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = c.SwitchToAbsorbing()
		}()
		go func() {
			defer wg.Done()
			c.SwitchToRestoring()
		}()
	}
	wg.Wait()

	// After the dust settles the absorber must be in a coherent state: one of
	// the valid machine states, never an intermediate torn value.
	endAbs := abs.State()
	switch endAbs {
	case model.AbsorberIdle, model.AbsorberAbsorbing:
		// valid terminal states
	default:
		t.Fatalf("absorber ended in invalid state %q after concurrent switches", endAbs)
	}

	// The controller state must also be a coherent voltage state, not a torn
	// read of a field that was written without the mutex.
	endState := c.State()
	switch endState {
	case model.VoltageAbsorbing, model.VoltageRestoring:
	default:
		t.Fatalf("controller ended in invalid state %q", endState)
	}

	// SwitchCount must equal the number of switch operations that actually
	// ran (2 * rounds). A torn counter would drift from this.
	if got := c.SwitchCount(); got != 2*rounds {
		t.Fatalf("switch count = %d, want %d", got, 2*rounds)
	}
}
