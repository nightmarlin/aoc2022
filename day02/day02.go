package day02

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/nightmarlin/aoc2022/lib"
)

type Day02 struct {
	log *zap.Logger
}

func New(log *zap.Logger) Day02 {
	return Day02{log: log}
}

// region Rock Paper Scissors

type Outcome int

const (
	LoseOutcome = iota - 1
	DrawOutcome
	WinOutcome
)

func (o Outcome) Score() int {
	return (int(o) + 1) * 3
}

type RPS interface {
	// VS returns the Outcome of comparing this RPS to another, from the
	// perspective of this RPS.
	VS(RPS) Outcome
	// TargetOutcome should return an RPS that, when VS is called, fulfills the
	// condition
	//
	//  rps.TargetOutcome(outcome).VS(rps) == outcome
	TargetOutcome(Outcome) RPS
	Value() int
}

type Rock string

func (Rock) VS(rps RPS) Outcome {
	switch rps.(type) {
	case Paper:
		return LoseOutcome
	case Scissors:
		return WinOutcome
	}
	return DrawOutcome
}
func (Rock) Value() int {
	return 1
}
func (Rock) TargetOutcome(o Outcome) RPS {
	switch o {
	case LoseOutcome:
		return Scissors("")
	case WinOutcome:
		return Paper("")
	}
	return Rock("")
}

type Paper string

func (Paper) VS(rps RPS) Outcome {
	switch rps.(type) {
	case Scissors:
		return LoseOutcome
	case Rock:
		return WinOutcome
	}
	return DrawOutcome
}
func (Paper) Value() int {
	return 2
}
func (Paper) TargetOutcome(o Outcome) RPS {
	switch o {
	case LoseOutcome:
		return Rock("")
	case WinOutcome:
		return Scissors("")
	}
	return Paper("")
}

type Scissors string

func (Scissors) VS(rps RPS) Outcome {
	switch rps.(type) {
	case Rock:
		return LoseOutcome
	case Paper:
		return WinOutcome
	}
	return DrawOutcome
}
func (Scissors) Value() int {
	return 3
}
func (Scissors) TargetOutcome(o Outcome) RPS {
	switch o {
	case LoseOutcome:
		return Paper("")
	case WinOutcome:
		return Rock("")
	}
	return Scissors("")
}

// endregion

func ToRPS(in uint8) RPS {
	switch in {
	case 'A', 'X':
		return Rock(in)
	case 'B', 'Y':
		return Paper(in)
	case 'C', 'Z':
		return Scissors(in)
	}
	return nil
}

func ToTargetOutcome(in uint8) Outcome {
	switch in {
	case 'X':
		return LoseOutcome
	case 'Z':
		return WinOutcome
	}
	return DrawOutcome
}

// Round runs a round, comparing your RPS and the opponent's RPS, and returning
// the score from your perspective (result + rps value)
func Round(theirs, yours RPS) int {
	return yours.VS(theirs).Score() + yours.Value()
}

// RunGame takes the puzzle input string, splits it into lines, and passes each
// line to the lineHandler to calculate the score. It then sums the resulting
// scores and returns that value.
func RunGame(input string, lineHandler func(theirCh, yourCh uint8) int) int {
	return lib.Reduce(
		lib.Map(
			strings.Split(input, "\n"),
			func(line string) int {
				if len(line) != 3 {
					return 0
				}
				return lineHandler(line[0], line[2])
			},
		),
		func(prev, next int) int { return prev + next },
		0,
	)
}

// PartOne presumes that the input is a guide to which RPS to play each round -
// If your opponent plays Rock then you should play Paper etc...
//
// The Mappings are defined as follows:
//
//	A, X : Rock     - Worth 1 point
//	B, Y : Paper    - Worth 2 points
//	C, Z : Scissors - Worth 3 points
//
// Your opponents move is first, denoted with ABC - followed by your move,
// denoted with XYZ.
//
// Your score is determined by the number of times you win, draw or lose (worth
// 6, 3 and 0 points respectively), plus the values of each RPS you played,
// defined above.
//
// Calculate the score from the given input using the above rules.
func (d Day02) PartOne(_ context.Context, input string) error {
	totalScore := RunGame(
		input,
		func(theirCh, yourCh uint8) (roundScore int) {
			return Round(ToRPS(theirCh), ToRPS(yourCh))
		},
	)

	d.log.Info("score calculated", zap.Int("score", totalScore))

	return nil
}

// PartTwo updates the definition to say that actually, XYZ refer to whether you
// should win, lose, or draw that round. Your opponent's move definitions, the
// value of each RPS, and the value of each Outcome remain the same.
//
//	X : You need to lose
//	Y : You need to draw
//	Z : You need to win
//
// Calculate the score from the given input using the above rules.
func (d Day02) PartTwo(_ context.Context, input string) error {
	totalScore := RunGame(
		input,
		func(theirCh, yourCh uint8) (roundScore int) {
			var (
				theirs      = ToRPS(theirCh)
				needOutcome = ToTargetOutcome(yourCh)
				yours       = theirs.TargetOutcome(needOutcome)

				val = Round(theirs, yours)
			)

			d.log.Debug(
				"round complete",
				zap.String("input", fmt.Sprintf("%c %c", theirCh, yourCh)),
				zap.String("theirs", fmt.Sprintf("%T", theirs)),
				zap.Any("want", needOutcome),
				zap.String("yours", fmt.Sprintf("%T", yours)),
				zap.Int("value", val),
			)

			return val
		},
	)

	d.log.Info("optimum score calculated", zap.Int("score", totalScore))

	return nil
}
