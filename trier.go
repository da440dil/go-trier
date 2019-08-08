// Package trier provides functions for (re-)execution of functions.
package trier

import (
	"context"
	"math/rand"
	"time"
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type trierError string

func (e trierError) Error() string {
	return string(e)
}

// ErrInvalidRetryCount is the error returned when WithCounter receives invalid value.
const ErrInvalidRetryCount = trierError("trier: number of retries must be greater than or equal to 0")

// ErrInvalidRetryDelay is the error returned when WithRetryDelay receives invalid value.
const ErrInvalidRetryDelay = trierError(
	"trier: delay between retries must be greater than or equal to 1 millisecond" +
		" and must be greater than or equal to jitter",
)

// ErrInvalidRetryJitter is the error returned when WithRetryJitter receives invalid value.
const ErrInvalidRetryJitter = trierError(
	"trier: retry jitter must be greater than or equal to 1 millisecond" +
		" and must be less than or equal to delay",
)

// ErrRetryCountExceeded is the error returned when number of retries exceeded.
const ErrRetryCountExceeded = trierError("trier: counter exceeded")

// Option is function returned by functions for setting options.
type Option func(*Trier) error

// WithRetryCount sets maximum number of retries.
// Must be greater than or equal to 0.
// By default equals 0.
func WithRetryCount(v int) Option {
	return func(t *Trier) error {
		if v < 0 {
			return ErrInvalidRetryCount
		}
		t.retryCount = v
		return nil
	}
}

// WithRetryDelay sets delay between retries.
// Must be greater than or equal to 1 millisecond.
// Must be greater than or equal to value of jitter.
// By default equals 0.
func WithRetryDelay(v time.Duration) Option {
	return func(t *Trier) error {
		if v < time.Millisecond {
			return ErrInvalidRetryDelay
		}
		t.retryDelay = durationToMilliseconds(v)
		if t.retryDelay < t.retryJitter {
			return ErrInvalidRetryDelay
		}
		return nil
	}
}

// WithRetryJitter sets maximum time randomly added to delay between retries
// to improve performance under high contention.
// Must be greater than or equal to 1 millisecond.
// Must be less than or equal to delay.
// By default equals 0.
func WithRetryJitter(v time.Duration) Option {
	return func(t *Trier) error {
		if v < time.Millisecond {
			return ErrInvalidRetryJitter
		}
		t.retryJitter = durationToMilliseconds(v)
		if t.retryJitter > t.retryDelay {
			return ErrInvalidRetryJitter
		}
		return nil
	}
}

// WithContext sets context which allows cancelling tries prematurely.
func WithContext(v context.Context) Option {
	return func(t *Trier) error {
		t.ctx = v
		return nil
	}
}

func durationToMilliseconds(d time.Duration) int {
	return int(d / time.Millisecond)
}

func millisecondsToDuration(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

// Trier defines parameters for executing retriable functions.
type Trier struct {
	retryCount  int // must be greater than or equal to zero
	retryDelay  int // must be greater than or equal to jitter
	retryJitter int // must be less than or equal to delay
	ctx         context.Context
}

// New creates new Trier.
func New(options ...Option) (*Trier, error) {
	r := &Trier{ctx: context.Background()}
	for _, fn := range options {
		if err := fn(r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// Creates new delay value based on initial delay value and initial value of jitter.
func (t *Trier) newDelay() int {
	if t.retryJitter == 0 {
		return t.retryDelay
	}
	min := t.retryDelay - t.retryJitter
	max := t.retryDelay + t.retryJitter
	v := rnd.Intn(max-min+1) + min
	return v
}

// Retriable is a function which execution could be retried.
type Retriable func() (bool, error)

// Try executes retriable function.
func (t *Trier) Try(fn Retriable) error {
	var counter = t.retryCount
	var timer *time.Timer
	for {
		ok, err := fn()
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
		if counter <= 0 {
			return ErrRetryCountExceeded
		}

		counter--
		timeout := millisecondsToDuration(t.newDelay())
		if timer == nil {
			timer = time.NewTimer(timeout)
			defer timer.Stop()
		} else {
			timer.Reset(timeout)
		}

		select {
		case <-t.ctx.Done():
			return t.ctx.Err()
		case <-timer.C:
		}
	}
}

// Try executes retriable function.
func Try(fn Retriable, options ...Option) error {
	r, err := New(options...)
	if err != nil {
		return err
	}
	return r.Try(fn)
}
