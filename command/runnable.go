package command

import (
	"context"
)

type Runnable interface {
	Run(ctx context.Context) error
}

type RunFunc func(ctx context.Context) error

func NewRunnable(unmarshalInterface UnmarshalInterface, f RunFunc) Runnable {
	return &runnableStruct{
		UnmarshalInterface: unmarshalInterface,
		f:                  f,
	}
}

type runnableStruct struct {
	UnmarshalInterface
	f RunFunc
}

func (r *runnableStruct) Run(ctx context.Context) error {
	return r.f(ctx)
}

type UnmarshalInterface interface{}
