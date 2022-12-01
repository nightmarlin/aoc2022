package day02

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type Day02 struct {
	log *zap.Logger
}

func New(log *zap.Logger) Day02 {
	return Day02{log: log}
}

func (d Day02) PartOne(ctx context.Context, input string) error {
	return errors.New("not implemented")
}

func (d Day02) PartTwo(ctx context.Context, input string) error {
	return errors.New("not implemented")
}
