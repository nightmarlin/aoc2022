package aoc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	RootURL    = "https://adventofcode.com/"
	URLPattern = "https://adventofcode.com/2022/day/%s/input"
)

type Fetcher struct {
	log         *zap.Logger
	client      *http.Client
	localFolder string
}

func NewFetcher(log *zap.Logger, sessionCookie string, localFolder string) (Fetcher, error) {
	cj, err := cookiejar.New(nil)
	if err != nil {
		return Fetcher{}, fmt.Errorf("failed to init cookie-jar: %w", err)
	}

	aocURL, err := url.Parse(RootURL)
	if err != nil {
		return Fetcher{}, fmt.Errorf("failed to parse root url: %w", err)
	}

	cj.SetCookies(
		aocURL,
		[]*http.Cookie{{Name: "session", Value: sessionCookie, MaxAge: 300}},
	)

	localFolder, err = filepath.Abs(localFolder)
	if err != nil {
		return Fetcher{}, fmt.Errorf("failed to find absolute path: %w", err)
	}

	err = os.MkdirAll(localFolder, os.ModePerm)
	if err != nil {
		return Fetcher{}, fmt.Errorf("failed to ensure local folder exists: %w", err)
	}

	return Fetcher{
			log: log.
				Named("input-fetcher").
				With(zap.String("localFolder", localFolder)),
			client:      &http.Client{Jar: cj},
			localFolder: localFolder,
		},
		nil
}

func (f Fetcher) FetchInput(ctx context.Context, day string) (string, error) {
	exists, err := f.isInputInLocalFolder(day)
	if err != nil {
		f.log.Warn("failed to check if input exists in local folder, fetching from aoc", zap.Error(err))

	} else if exists {
		f.log.Info("input found in local folder, will load from there")

		input, err := f.fetchInputFromLocalFolder(day)
		if err != nil {
			f.log.Warn("failed to load input from local folder, fetching from aoc", zap.Error(err))
		} else {
			return input, nil
		}

	} else {
		f.log.Info("input not found in local folder, fetching from aoc")
	}

	input, err := f.fetchInputFromAOC(ctx, day)
	if err != nil {
		return "", fmt.Errorf("failed to fetch input from aoc: %w", err)
	}

	f.log.Info("fetched input from aoc")

	err = f.saveInputToLocalFolder(day, input)
	if err != nil {
		f.log.Warn("failed to save input to local folder", zap.Error(err))
	} else {
		f.log.Info("saved input to local folder, future runs will use ths version")
	}

	return input, nil
}

func (f Fetcher) fetchInputFromAOC(ctx context.Context, day string) (string, error) {
	u := fmt.Sprintf(URLPattern, day)
	f.log.Debug("fetching input", zap.String("url", u))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create http request: %w", err)
	}

	res, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform http request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		body = []byte(fmt.Sprintf("<failed to read request body: %s>", err.Error()))
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got status %d, wanted 200. body: %s", res.StatusCode, body)
	}

	return string(body), nil
}

// region filesystem

func (f Fetcher) inputFileName(day string) string {
	if len(day) == 1 {
		day = fmt.Sprintf("0%s", day)
	}
	return filepath.Join(f.localFolder, day)
}

func (f Fetcher) isInputInLocalFolder(day string) (bool, error) {
	_, err := os.Stat(f.inputFileName(day))
	switch {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("unable to check if input for day exists: %w", err)
	default:
		return true, nil
	}
}

func (f Fetcher) fetchInputFromLocalFolder(day string) (string, error) {
	input, err := os.ReadFile(f.inputFileName(day))
	if err != nil {
		return "", fmt.Errorf("failed to read input file for day: %w", err)
	}
	return string(input), nil
}

func (f Fetcher) saveInputToLocalFolder(day, input string) error {
	err := os.WriteFile(f.inputFileName(day), []byte(input), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write input file for day: %w", err)
	}
	return nil
}

// endregion
