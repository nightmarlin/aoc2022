package day01

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/nightmarlin/aoc2022/lib"
)

// Day01 is a challenge focused on basic string parsing. It can be solved using
// a set of map/reduce operations.
type Day01 struct {
	log *zap.Logger
}

func New(log *zap.Logger) Day01 {
	return Day01{log: log.Named("day-01")}
}

// ParseLineValues takes a single line of the input string and parses it into
// a whole number.
func (d Day01) ParseLineValues(valueStr string) int {
	if valueStr == "" {
		return 0
	}

	val, err := strconv.ParseInt(valueStr, 10, 32) // attempt to parse line
	if err != nil {
		d.log.Warn("failed to parse line", zap.String("entry", valueStr))
		val = 0
	}
	return int(val)
}

// SumEachGroup splits the input string into groups (two consecutive newlines),
// maps the lines in each group to a number (using ParseLineValues), and then
// sums (reduces) the lines in each group to a single value.
func (d Day01) SumEachGroup(input string) []int {
	return lib.Map(
		strings.Split(input, "\n\n"), // Each elf is split by two lines
		func(group string) int {
			return lib.Reduce(
				lib.Map(
					strings.Split(group, "\n"), // Each calorie count is on its own line
					d.ParseLineValues,
				),
				func(prev, next int) int { return prev + next }, // Sum total calories for each elf
				0,
			)
		},
	)
}

// SortDesc sorts the input slice from high to low.
func (d Day01) SortDesc(groupSums []int) []int {
	sort.Slice(groupSums, func(i, j int) bool { return groupSums[i] > groupSums[j] })
	return groupSums
}

// SumTopN slices the first n values of a sorted slice and sums them together.
func (d Day01) SumTopN(sortedGroupSums []int, n int) int {
	if len(sortedGroupSums) == 0 {
		return 0
	}

	if len(sortedGroupSums) < n {
		n = len(sortedGroupSums)
	}

	return lib.Reduce(
		sortedGroupSums[:n],
		func(prev int, next int) int { return prev + next },
		0,
	)
}

// PartOne asks to find the highest number of calories held by a single Elf.
// The input is a newline-delimited list of integers of the form:
//
//	123
//	456
//
//	789
//	12
//
// Where each line with a value represents the calorie count for an item held by
// that Elf, and each grouping of items represents the set of items held by that
// Elf.
func (d Day01) PartOne(_ context.Context, input string) error {
	mostCalories := d.SumTopN(d.SortDesc(d.SumEachGroup(input)), 1)

	d.log.Info(
		"maximum calorie count found",
		zap.Int("calories", mostCalories),
	)
	return nil
}

// PartTwo asks a similar question, but in the spirit of fairness asks the total
// number of calories shared between the three Elves carrying the most calories.
func (d Day01) PartTwo(_ context.Context, input string) error {
	topThreeSum := d.SumTopN(d.SortDesc(d.SumEachGroup(input)), 3)

	d.log.Info(
		"sum of calories for 3 elves holding most calories found",
		zap.Int("calories-sum", topThreeSum),
	)
	return nil
}
