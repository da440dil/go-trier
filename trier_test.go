package trier

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type imock struct {
	d    time.Duration
	done bool
}

func (m *imock) Next() (time.Duration, bool) {
	return m.d, m.done
}

func (m *imock) Iterator() Iterator {
	return m
}

type trymock struct {
	ok  bool
	err error
}

func (m *trymock) Try(ctx context.Context) (bool, error) {
	return m.ok, m.err
}

func TestTrier(t *testing.T) {
	b := &imock{0, true}
	tr := Trier{b}

	e := errors.New("some error")
	f := &trymock{ok: false, err: e}
	ctx := context.Background()

	ok, err := tr.Try(ctx, f.Try)
	require.Equal(t, e, err)
	require.False(t, ok)

	f.ok = true
	f.err = nil
	ok, err = tr.Try(ctx, f.Try)
	require.NoError(t, err)
	require.True(t, ok)

	f.ok = false
	ok, err = tr.Try(ctx, f.Try)
	require.NoError(t, err)
	require.False(t, ok)

	b.d = time.Millisecond * 100
	b.done = false
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*150)
	defer cancel()
	ok, err = tr.Try(ctx, f.Try)
	require.Equal(t, context.DeadlineExceeded, err)
	require.False(t, ok)
}

func TestNewTrier(t *testing.T) {
	b := &imock{0, true}
	w := &imock{0, false}
	tr := NewTrier(b, func(Iterable) Iterable {
		return w
	})
	require.Equal(t, w, tr.b)
}
