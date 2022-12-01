package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/nightmarlin/aoc2022/aoc"
	"github.com/nightmarlin/aoc2022/day01"
	"github.com/nightmarlin/aoc2022/lib"
)

type Challenge interface {
	PartOne(ctx context.Context, input string) error
	PartTwo(ctx context.Context, input string) error
}

var log *zap.Logger

var challenges = map[string]func() Challenge{
	"1": func() Challenge { return day01.New(log) },
}

func init() {
	l, err := zap.NewDevelopment()
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "failed to init logger:", err)
		os.Exit(1)
	}
	log = l
}

func main() {
	ctx := context.Background()

	sessionCookie := os.Getenv("SESSION_COOKIE")
	if sessionCookie == "" {
		log.Fatal("session cookie must be set to fetch aoc inputs")
	}

	log.Info(
		"please choose a challenge to run",
		zap.Strings("available days", lib.Keys(challenges)),
	)

	day, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal("failed to read input from stdin", zap.Error(err))
	}

	day = strings.Trim(day, "0\n")

	cInit, ok := challenges[day]
	if !ok {
		log.Fatal(
			"the specified challenge has not been completed or does not exist",
			zap.String("day", day),
		)
	}

	fetcher, err := aoc.NewFetcher(log, sessionCookie, "inputs")
	if err != nil {
		log.Fatal("failed to init aoc fetcher", zap.Error(err))
	}

	input, err := fetcher.FetchInput(ctx, day)
	if err != nil {
		log.Fatal("unable to get input for chosen day", zap.Error(err))
	}

	log.Info("input fetched, initializing solution")
	challenge := cInit()

	log.Info("solution initialized, running part one...")
	err = challenge.PartOne(ctx, input)
	if err != nil {
		log.Fatal("error occurred while running solution part one", zap.Error(err))
	}

	log.Info("solution initialized, running part two...")
	err = challenge.PartTwo(ctx, input)
	if err != nil {
		log.Fatal("error occurred while running solution part two", zap.Error(err))
	}

	log.Info("complete!")
}
