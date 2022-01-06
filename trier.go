// Package trier provides data structures for re-execution of functions within configurable limits.
package trier

import (
	"context"
	"time"
)

// Trier defines parameters for executing retriable functions.
type Trier struct {
	b Iterable
}

// NewTrier creates new trier.
func NewTrier(b Iterable, fns ...Decorator) Trier {
	for _, fn := range fns {
		b = fn(b)
	}
	return Trier{b}
}

// Retriable is a function which execution could be retried, returns execution success flag.
type Retriable func(ctx context.Context) (bool, error)

// Try executes retriable function, retries execution if execution success flag equals false.
func (t Trier) Try(ctx context.Context, fn Retriable) (bool, error) {
	var it Iterator
	var timer *time.Timer
	for {
		ok, err := fn(ctx)
		if err != nil {
			return false, err
		}
		if ok {
			return ok, err
		}
		if it == nil {
			it = t.b.Iterator()
		}
		d, done := it.Next()
		if done {
			return false, err
		}
		if timer == nil {
			timer = time.NewTimer(d)
			defer timer.Stop()
		} else {
			timer.Reset(d)
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-timer.C:
		}
	}
}
