package trier

import (
	"math/rand"
	"time"
)

// Iterator defines parameters to create new delay.
type Iterator interface {
	Next() (time.Duration, bool)
}

// Iterable defines parameters to create new iterator.
type Iterable interface {
	Iterator() Iterator
}

type constant time.Duration

func (i constant) Next() (time.Duration, bool) {
	return time.Duration(i), false
}
func (i constant) Iterator() Iterator {
	return i
}

// NewConstant creates iterable for creating iterators
// returning constant delay on each iteration.
func NewConstant(d time.Duration) Iterable {
	return constant(d)
}

type linear struct {
	d, rate time.Duration
}

func (i *linear) Next() (time.Duration, bool) {
	i.d += i.rate
	return i.d, false
}

func (i linear) Iterator() Iterator {
	return &linear{rate: i.d}
}

// NewLinear creates iterable for creating iterators
// returning delay which grows linearly on each iteration.
func NewLinear(d time.Duration) Iterable {
	return linear{d: d}
}

type linearRate struct {
	d, rate time.Duration
}

func (i *linearRate) Next() (time.Duration, bool) {
	v := i.d
	i.d += i.rate
	return v, false
}

func (i linearRate) Iterator() Iterator {
	return &linearRate{i.d, i.rate}
}

// NewLinearRate creates iterable for creating iterators
// returning delay which grows linearly with specified rate on each iteration.
func NewLinearRate(d, rate time.Duration) Iterable {
	return linearRate{d, rate}
}

type exponential time.Duration

func (i *exponential) Next() (time.Duration, bool) {
	v := *i
	*i = v + v
	return time.Duration(v), false
}

func (i exponential) Iterator() Iterator {
	return &i
}

// NewExponential creates iterable for creating iterators
// returning delay which grows exponentially on each iteration.
func NewExponential(d time.Duration) Iterable {
	return exponential(d)
}

type exponentialRate struct {
	d, rate float64
}

func (i *exponentialRate) Next() (time.Duration, bool) {
	v := i.d
	i.d += i.d * i.rate
	return time.Duration(v), false
}

func (i exponentialRate) Iterator() Iterator {
	return &exponentialRate{i.d, i.rate}
}

// NewExponentialRate creates iterable for creating iterators
// returning delay which grows exponentially with specified rate on each iteration.
func NewExponentialRate(d time.Duration, rate float64) Iterable {
	return exponentialRate{float64(d), rate}
}

// Decorator extends behavior of an iterable.
type Decorator func(Iterable) Iterable

type maxRetriesB struct {
	n int
	b Iterable
}

func (b maxRetriesB) Iterator() Iterator {
	return &maxRetriesI{b.n, b.b.Iterator()}
}

type maxRetriesI struct {
	n int
	i Iterator
}

func (i *maxRetriesI) Next() (time.Duration, bool) {
	if i.n > 0 {
		i.n--
		return i.i.Next()
	}
	return 0, true
}

// WithMaxRetries sets maximum number of retries.
func WithMaxRetries(n int) Decorator {
	return func(b Iterable) Iterable {
		return maxRetriesB{n, b}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type jitterB struct {
	n, j int64
	b    Iterable
}

func (b jitterB) Iterator() Iterator {
	return jitterI{b.n, b.j, b.b.Iterator()}
}

type jitterI struct {
	n, j int64
	i    Iterator
}

func (i jitterI) Next() (time.Duration, bool) {
	v, done := i.i.Next()
	if done {
		return 0, done
	}
	v = v + time.Duration(rand.Int63n(i.n)-i.j)
	if v < 0 {
		v = 0
	}
	return v, done
}

// WithJitter sets maximum duration randomly added to or extracted from delay
// between retries to improve performance under high contention.
func WithJitter(d time.Duration) Decorator {
	return func(b Iterable) Iterable {
		j := int64(d)
		return jitterB{j*2 + 1, j, b}
	}
}
