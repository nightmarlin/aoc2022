package day03

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"github.com/nightmarlin/aoc2022/lib"
)

type Day03 struct {
	log *zap.Logger
}

func New(log *zap.Logger) Day03 {
	return Day03{log: log.Named("day-03")}
}

// region set

type Set[T comparable] interface {
	Members() []T            // Members returns the items in the set in a slice.
	Insert(T)                // Insert adds an item to the set. If it exists already, nothing changes.
	Intersect(Set[T]) Set[T] // Intersect finds the Set of items that exist in both Sets.
}

// MapSet implements a Set based on a map. It is not safe for concurrent use.
type MapSet[T comparable] map[T]struct{}

// NewMapSet initializes and returns an empty MapSet implementation of Set.
func NewMapSet[T comparable]() MapSet[T] { return make(MapSet[T]) }

func (m MapSet[T]) Members() []T {
	res := make([]T, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
func (m MapSet[T]) Insert(t T) { m[t] = struct{}{} }
func (m MapSet[T]) Intersect(other Set[T]) Set[T] {
	mMembers, oMembers := m.Members(), other.Members()
	intersection := lib.Filter(
		mMembers,
		func(mMember T) bool {
			return lib.Any(
				oMembers,
				func(oMember T) bool { return oMember == mMember },
			)
		},
	)

	resSet := NewMapSet[T]()
	for i := range intersection {
		resSet.Insert(intersection[i])
	}
	return resSet
}

// endregion

// ConstructSet creates a new Set from the characters in the input string.
func ConstructSet(chars string) Set[uint8] {
	res := NewMapSet[uint8]()
	for i := range chars {
		res.Insert(chars[i])
	}
	return res
}

// LineToCompartments splits the input string in half and returns a Set for each
// half of the string using ConstructSet.
func LineToCompartments(line string) (Set[uint8], Set[uint8]) {
	halfLineLen := len(line) / 2
	return ConstructSet(line[:halfLineLen]), ConstructSet(line[halfLineLen:])
}

// The Priority for lowercase characters is 1-26 (a-z), and for uppercase
// characters it 27-52 (A-Z). Any other character is undefined.
func Priority(c uint8) int {
	switch {
	case 'a' <= c && c <= 'z':
		return int(c-'a') + 1
	case 'A' <= c && c <= 'Z':
		return int(c-'A') + 27
	}
	return 0
}

// PartOne asks us to find the item that exists in both compartments of a bag.
// Each bag is represented by a single line of input, with each half of the
// string representing a compartment in the bag.
//
// We are then asked to assign this item a Priority and return the sum of these
// priorities for every line of the input.
//
// The priority is defined as:
//
//	a-z: 1-26 inclusive
//	A-Z: 27-52 inclusive
//
// This solution models each compartment as a set and attempts to find the
// intersection
func (d Day03) PartOne(_ context.Context, input string) error {
	prioritySum := lib.Reduce(
		lib.Map(
			strings.Split(input, "\n"),
			func(line string) uint8 {
				if len(line) == 0 {
					return 0
				}
				compartment1, compartment2 := LineToCompartments(line)      // Convert each line to the two compartments
				intersect := compartment1.Intersect(compartment2).Members() // Find the intersection
				if len(intersect) == 0 {
					return 0
				}
				return intersect[0]
			},
		),
		func(prev int, next uint8) int { return prev + Priority(next) }, // Calculate and sum the priority for each element
		0,
	)

	d.log.Info(
		"found the sum of the priorities for items in both compartments of each bag",
		zap.Int("sum", prioritySum),
	)
	return nil
}

// PartTwo builds upon PartOne by looking across multiple bags at once. For each
// set of three consecutive bags, there should be a single item that is held in
// all three bags. We are now asked to identify this item, assign it a Priority
// as in PartOne and return the sum of the priorities across every three-bag
// group.
func (d Day03) PartTwo(_ context.Context, input string) error {
	var (
		total        int
		intersectSet Set[uint8]
	)
	for i, line := range strings.Split(input, "\n") {
		// On the first iteration, and every three iterations after...
		if i%3 == 0 {
			// If the set exists
			if intersectSet != nil {
				// Fetch the set of items that are in all three bags and add the
				// priority of the first to the running total. There should only be one
				// item.
				members := intersectSet.Members()
				if len(members) != 0 {
					total += Priority(members[0])
				}
			}

			intersectSet = ConstructSet(line)
		} else {

			// On any other iteration, convert each bag to a set and find the
			// intersection
			intersectSet = intersectSet.Intersect(ConstructSet(line))
		}
	}

	d.log.Info(
		"found the sum of the priorities for items in the bags of every elf in each group",
		zap.Int("sum", total),
	)
	return nil
}
