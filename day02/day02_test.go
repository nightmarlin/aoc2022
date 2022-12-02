package day02

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTargetOutcomes(t *testing.T) {
	testTable := []struct {
		Name string

		RPS     RPS
		Outcome Outcome
	}{
		{Name: "win against rock", RPS: Rock(""), Outcome: WinOutcome},
		{Name: "draw against rock", RPS: Rock(""), Outcome: DrawOutcome},
		{Name: "loss against rock", RPS: Rock(""), Outcome: LoseOutcome},
		{Name: "win against paper", RPS: Paper(""), Outcome: WinOutcome},
		{Name: "draw against paper", RPS: Paper(""), Outcome: DrawOutcome},
		{Name: "loss against paper", RPS: Paper(""), Outcome: LoseOutcome},
		{Name: "win against scissors", RPS: Scissors(""), Outcome: WinOutcome},
		{Name: "draw against scissors", RPS: Scissors(""), Outcome: DrawOutcome},
		{Name: "loss against scissors", RPS: Scissors(""), Outcome: LoseOutcome},
	}

	for _, entry := range testTable {
		entry := entry
		t.Run(
			entry.Name,
			func(t *testing.T) {
				t.Parallel()

				assert.Equal(
					t,
					entry.RPS.TargetOutcome(entry.Outcome).VS(entry.RPS),
					entry.Outcome,
				)
			},
		)
	}
}
