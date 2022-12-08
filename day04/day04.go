package day04

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/nightmarlin/aoc2022/lib"
)

type Day04 struct {
	log *zap.Logger
}

func New(log *zap.Logger) Day04 {
	return Day04{log: log.Named("day-04")}
}

type Range [2]int

func (r Range) Min() int { return r[0] }
func (r Range) Max() int { return r[1] }

func (r Range) Contains(o Range) bool {
	return r.Min() <= o.Min() && o.Max() <= r.Max()
}

func (r Range) Intersects(o Range) bool {
	return r.Min() <= o.Max() && r.Max() >= o.Min()
}

type Row [2]Range

func (r Row) EitherContains() bool {
	return r[0].Contains(r[1]) || r[1].Contains(r[0])
}

func (r Row) Intersect() bool {
	return r[0].Intersects(r[1])
}

func GetRange(r string) (Range, error) {
	bounds := strings.Split(r, "-")
	if len(bounds) != 2 {
		return Range{}, fmt.Errorf("a range should contain an upper and lower bound, got: %q", r)
	}

	var (
		res Range
		err error
	)

	res[0], err = strconv.Atoi(bounds[0])
	if err != nil {
		return Range{}, fmt.Errorf("failed to parse lower bound: %w", err)
	}
	res[1], err = strconv.Atoi(bounds[1])
	if err != nil {
		return Range{}, fmt.Errorf("failed to parse upper bound: %w", err)
	}
	return res, nil
}

func GetRow(line string) (Row, error) {
	ranges := strings.Split(line, ",")
	if len(ranges) != 2 {
		return Row{}, fmt.Errorf("a row should contain 2 ranges, got %d: %q", len(ranges), line)
	}

	var (
		res Row
		err error
	)

	res[0], err = GetRange(ranges[0])
	if err != nil {
		return Row{}, fmt.Errorf("failed to parse range: %w", err)
	}
	res[1], err = GetRange(ranges[1])
	if err != nil {
		return Row{}, fmt.Errorf("failed to parse range: %w", err)
	}
	return res, nil
}

func (d Day04) PartOne(_ context.Context, input string) error {
	containCount := lib.Reduce(
		lib.Map(
			lib.Filter(strings.Split(input, "\n"), func(s string) bool { return len(s) > 0 }),
			func(line string) bool {
				r, err := GetRow(line)
				if err != nil {
					d.log.Warn("failed to parse row", zap.Error(err))
					return false
				}
				ec := r.EitherContains()

				d.log.Debug("row parsed", zap.Any("row", r), zap.Bool("either_contains", ec))
				return ec
			},
		),
		func(sum int, b bool) int {
			if !b {
				return sum
			}
			return sum + 1
		},
		0,
	)

	d.log.Info("found number of pairs where one fully contains the other", zap.Int("count", containCount))
	return nil
}

func (d Day04) PartTwo(_ context.Context, input string) error {
	intersectCount := lib.Reduce(
		lib.Map(
			lib.Filter(strings.Split(input, "\n"), func(s string) bool { return len(s) > 0 }),
			func(line string) bool {
				r, err := GetRow(line)
				if err != nil {
					d.log.Warn("failed to parse row", zap.Error(err))
					return false
				}
				ec := r.Intersect()

				d.log.Debug("row parsed", zap.Any("row", r), zap.Bool("intersects", ec))
				return ec
			},
		),
		func(sum int, b bool) int {
			if !b {
				return sum
			}
			return sum + 1
		},
		0,
	)

	d.log.Info("found number of pairs where one intersects with the other", zap.Int("count", intersectCount))
	return nil
}
