package trier

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	t.Run("counter < 0", func(t *testing.T) {
		_, err := New(WithRetryCount(-1))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryCount, err)
	})

	t.Run("delay < 1 millisecond", func(t *testing.T) {
		_, err := New(WithRetryDelay(time.Microsecond))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryDelay, err)
	})

	t.Run("delay < jitter", func(t *testing.T) {
		_, err := New(
			WithRetryDelay(time.Millisecond*3),
			WithRetryJitter(time.Millisecond*2),
			WithRetryDelay(time.Millisecond*1),
		)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryDelay, err)
	})

	t.Run("jitter < 1 millisecond", func(t *testing.T) {
		_, err := New(WithRetryJitter(time.Microsecond))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryJitter, err)
	})

	t.Run("jitter > delay", func(t *testing.T) {
		_, err := New(
			WithRetryDelay(time.Millisecond*2),
			WithRetryJitter(time.Millisecond*3),
		)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryJitter, err)
	})

	t.Run("success", func(t *testing.T) {
		rt, err := New(
			WithRetryCount(42),
			WithRetryDelay(time.Millisecond*100),
			WithRetryJitter(time.Millisecond*20),
		)
		assert.NoError(t, err)
		assert.IsType(t, &Trier{}, rt)
		assert.Equal(t, 42, rt.retryCount)
		assert.Equal(t, 100, rt.retryDelay)
		assert.Equal(t, 20, rt.retryJitter)
	})
}

func TestTrierNewDelay(t *testing.T) {
	t.Run("jitter = 0", func(t *testing.T) {
		delay := 42
		r := &Trier{retryDelay: delay, retryJitter: 0}
		v := r.newDelay()
		assert.Equal(t, delay, v)
	})

	testCases := []struct {
		delay  int
		jitter int
	}{
		{100, 20},
		{200, 50},
		{1000, 100},
		{500, 500},
	}

	for _, tc := range testCases {
		delay := tc.delay
		jitter := tc.jitter

		t.Run(fmt.Sprintf("delay = %v; jitter = %v", delay, jitter), func(t *testing.T) {
			r := &Trier{retryDelay: delay, retryJitter: jitter}
			v := r.newDelay()
			assert.True(t, v >= (delay-jitter) && v <= (delay+jitter))
		})
	}
}

type mock struct {
	ok      bool
	d       time.Duration
	err     error
	counter int
}

func (m *mock) Try() (bool, time.Duration, error) {
	m.counter++
	return m.ok, m.d, m.err
}

func TestNewTrier(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		rt, err := New()
		assert.NoError(t, err)

		e := errors.New("any")
		m := &mock{ok: false, err: e}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, e, err)
		assert.Equal(t, 1, m.counter)
	})

	t.Run("ok", func(t *testing.T) {
		rt, err := New()
		assert.NoError(t, err)

		m := &mock{ok: true, err: nil}
		err = rt.Try(m.Try)
		assert.NoError(t, err)
		assert.Equal(t, 1, m.counter)
	})

	t.Run("ErrCounterExceeded", func(t *testing.T) {
		rt, err := New()
		assert.NoError(t, err)

		m := &mock{ok: false, err: nil}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, ErrRetryCountExceeded, err)
		assert.Equal(t, 1, m.counter)
	})

	t.Run("ErrCounterExceeded WithCounter", func(t *testing.T) {
		rt, err := New(WithRetryCount(2))
		assert.NoError(t, err)

		m := &mock{ok: false, err: nil}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, ErrRetryCountExceeded, err)
		assert.Equal(t, 3, m.counter)
	})

	t.Run("ErrCounterExceeded WithCounter WithContext", func(t *testing.T) {
		rt, err := New(WithRetryCount(2), WithContext(context.Background()))
		assert.NoError(t, err)

		m := &mock{ok: false, err: nil}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, ErrRetryCountExceeded, err)
		assert.Equal(t, 3, m.counter)
	})

	t.Run("ErrCounterExceeded WithCounter WithDelay WithContext", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()

		rt, err := New(WithRetryCount(2), WithRetryDelay(time.Millisecond*50), WithContext(ctx))
		assert.NoError(t, err)

		m := &mock{ok: false, d: -1, err: nil}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, ErrRetryCountExceeded, err)
		assert.Equal(t, 3, m.counter)
	})

	t.Run("context.DeadlineExceeded WithCounter WithDelay WithContext", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		rt, err := New(WithRetryCount(2), WithRetryDelay(time.Millisecond*200), WithContext(ctx))
		assert.NoError(t, err)

		m := &mock{ok: false, d: -1, err: nil}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
		assert.Equal(t, 1, m.counter)
	})

	t.Run("ErrCounterExceeded with custom delay", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()

		rt, err := New(WithRetryCount(2), WithRetryDelay(time.Millisecond*200), WithContext(ctx))
		assert.NoError(t, err)

		m := &mock{ok: false, d: 0, err: nil}
		err = rt.Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, ErrRetryCountExceeded, err)
		assert.Equal(t, 3, m.counter)
	})
}

func TestTry(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		e := errors.New("any")
		m := &mock{ok: false, err: e}
		err := Try(m.Try)
		assert.Error(t, err)
		assert.Equal(t, e, err)
		assert.Equal(t, 1, m.counter)
	})

	t.Run("ok", func(t *testing.T) {
		m := &mock{ok: true, err: nil}
		err := Try(m.Try)
		assert.NoError(t, err)
		assert.Equal(t, 1, m.counter)
	})

	t.Run("counter < 0", func(t *testing.T) {
		err := Try(new(mock).Try, WithRetryCount(-1))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryCount, err)
	})

	t.Run("delay < 1 millisecond", func(t *testing.T) {
		err := Try(new(mock).Try, WithRetryDelay(time.Microsecond))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryDelay, err)
	})

	t.Run("delay < jitter", func(t *testing.T) {
		err := Try(
			new(mock).Try,
			WithRetryDelay(time.Millisecond*3),
			WithRetryJitter(time.Millisecond*2),
			WithRetryDelay(time.Millisecond*1),
		)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryDelay, err)
	})

	t.Run("jitter < 1 millisecond", func(t *testing.T) {
		err := Try(new(mock).Try, WithRetryJitter(time.Microsecond))
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryJitter, err)
	})

	t.Run("jitter > delay", func(t *testing.T) {
		err := Try(
			new(mock).Try,
			WithRetryDelay(time.Millisecond*2),
			WithRetryJitter(time.Millisecond*3),
		)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidRetryJitter, err)
	})
}

func TestTrierError(t *testing.T) {
	v := "any"
	err := trierError(v)
	assert.Equal(t, v, err.Error())
}
