package trier

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrierError(t *testing.T) {
	v := "any"
	err := trierError(v)
	assert.Equal(t, v, err.Error())
}

func TestTTLError(t *testing.T) {
	d := time.Millisecond * 42
	err := newTTLError(d)
	assert.True(t, errors.Is(err, ErrTooManyRetries))
	assert.Equal(t, ErrTooManyRetries.Error(), err.Error())
	assert.Equal(t, d, err.TTL())
}
