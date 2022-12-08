package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/nightmarlin/aoc2022/aoc"
	"github.com/nightmarlin/aoc2022/day01"
	"github.com/nightmarlin/aoc2022/day02"
	"github.com/nightmarlin/aoc2022/day03"
	"github.com/nightmarlin/aoc2022/day04"
	"github.com/nightmarlin/aoc2022/lib"
)

// A Solution carries out some task on the input, as defined by the AoC
// challenge. Solutions are expected to gracefully handle cancelled contexts if
// there is a likelihood of them running for extended time periods.
type Solution interface {
	// PartOne of a Solution is generally a specific application of a problem
	// statement.
	PartOne(ctx context.Context, input string) error

	// PartTwo of a Solution is typically a more generalised form of the problem
	// presented in PartOne, using the same input.
	PartTwo(ctx context.Context, input string) error
}

var solutions = map[string]func(*zap.Logger) Solution{
	"1": func(log *zap.Logger) Solution { return day01.New(log) },
	"2": func(log *zap.Logger) Solution { return day02.New(log) },
	"3": func(log *zap.Logger) Solution { return day03.New(log) },
	"4": func(log *zap.Logger) Solution { return day04.New(log) },
}

func initLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level.SetLevel(zap.InfoLevel)
	cfg.DisableStacktrace = true

	if os.Getenv("TRACE") != "" {
		cfg.Level.SetLevel(zap.DebugLevel)
		cfg.DisableStacktrace = false
	}

	log, err := cfg.Build()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to init logger: %s\n", err.Error())
		os.Exit(1)
	}

	return log
}

func main() {
	ctx := context.Background()
	log := initLogger()

	sessionCookie := os.Getenv("SESSION_COOKIE")
	if sessionCookie == "" {
		log.Fatal("session cookie must be set to fetch aoc inputs")
	}
	localFolder := os.Getenv("LOCAL_FOLDER")
	if localFolder == "" {
		localFolder = "inputs"
	}

	day := os.Getenv("SOLUTION")
	if day == "" {
		log.Info(
			"please choose a solution to run",
			zap.Ints("available days", sortedSolutionNames()),
		)

		read, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatal("failed to read input from stdin", zap.Error(err))
		}
		day = read
	} else {
		log.Info("solution set via environment variable", zap.String("day", day))
	}
	day = strings.Trim(day, "0\n")

	sInit, ok := solutions[day]
	if !ok {
		log.Fatal("the specified solution has not been completed or does not exist", zap.String("day", day))
	}

	fetcher, err := aoc.NewFetcher(log, sessionCookie, localFolder)
	if err != nil {
		log.Fatal("failed to init aoc fetcher", zap.Error(err))
	}

	input, err := fetcher.FetchInput(ctx, day)
	if err != nil {
		log.Fatal("unable to get input for chosen day", zap.Error(err))
	}

	log.Info("input fetched, initializing solution")
	solution := sInit(log)

	log.Info("solution initialized, running part one...")
	err = solution.PartOne(ctx, input)
	if err != nil {
		log.Fatal("error occurred while running solution part one", zap.Error(err))
	}

	log.Info("part one complete, running part two...")
	err = solution.PartTwo(ctx, input)
	if err != nil {
		log.Fatal("error occurred while running solution part two", zap.Error(err))
	}

	log.Info("complete!")
}

// Map iteration is non-deterministic, so we need to manually sort the keys.
func sortedSolutionNames() []int {
	asNums := lib.Map(
		lib.Keys(solutions),
		func(k string) int {
			i, err := strconv.ParseInt(k, 10, 32)
			if err != nil {
				return 0
			}
			return int(i)
		},
	)

	sort.SliceStable(asNums, func(i, j int) bool { return asNums[i] < asNums[j] })

	return asNums
}
