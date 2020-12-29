package trier

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewConstant(t *testing.T) {
	b := NewConstant(time.Second)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
	}
}

func ExampleNewConstant() {
	it := NewConstant(time.Second).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 1s, false }
	// #2: { 1s, false }
	// #3: { 1s, false }
}

func TestNewLinear(t *testing.T) {
	b := NewLinear(time.Second)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*2, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*3, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*4, d)
		require.False(t, done)
	}
}

func ExampleNewLinear() {
	it := NewLinear(time.Second).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 3s, false }
	// #3: { 4s, false }
}

func TestNewLinearRate(t *testing.T) {
	b := NewLinearRate(time.Second, time.Second*2)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*3, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*5, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*7, d)
		require.False(t, done)
	}
}

func ExampleNewLinearRate() {
	it := NewLinearRate(time.Second, time.Second*2).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 3s, false }
	// #2: { 5s, false }
	// #3: { 7s, false }
}

func TestNewExponential(t *testing.T) {
	b := NewExponential(time.Second)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*2, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*4, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second*8, d)
		require.False(t, done)
	}
}

func ExampleNewExponential() {
	it := NewExponential(time.Second).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 4s, false }
	// #3: { 8s, false }
}

func TestNewExponentialRate(t *testing.T) {
	b := NewExponentialRate(time.Second, 0.2)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Millisecond*1000, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*1200, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*1440, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Millisecond*1728, d)
		require.False(t, done)
	}
}

func ExampleNewExponentialRate() {
	it := NewExponentialRate(time.Second, 0.2).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 1.2s, false }
	// #2: { 1.44s, false }
	// #3: { 1.728s, false }
}

func TestWithMaxRetries(t *testing.T) {
	b := NewConstant(time.Second)
	b = WithMaxRetries(3)(b)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Second, d)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Duration(0), d)
		require.True(t, done)
	}
}

func ExampleWithMaxRetries_newConstant() {
	it := WithMaxRetries(3)(NewConstant(time.Second)).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 1s, false }
	// #2: { 1s, false }
	// #3: { 0s, true }
}

func ExampleWithMaxRetries_newLinear() {
	it := WithMaxRetries(3)(NewLinear(time.Second)).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 3s, false }
	// #3: { 0s, true }
}

func ExampleWithMaxRetries_newExponential() {
	it := WithMaxRetries(3)(NewExponential(time.Second)).Iterator()
	for i := 0; i < 4; i++ {
		d, done := it.Next()
		fmt.Printf("#%v: { %v, %v }\n", i, d, done)
	}
	// Output:
	// #0: { 1s, false }
	// #1: { 2s, false }
	// #2: { 4s, false }
	// #3: { 0s, true }
}

func TestWithJitter(t *testing.T) {
	b := NewLinear(time.Second)
	b = WithMaxRetries(3)(b)
	b = WithJitter(time.Millisecond * 100)(b)
	for i := 0; i < 3; i++ {
		it := b.Iterator()
		d, done := it.Next()
		require.True(t, time.Millisecond*900 <= d && d <= time.Millisecond*1100)
		require.False(t, done)
		d, done = it.Next()
		require.True(t, time.Millisecond*1900 <= d && d <= time.Millisecond*2100)
		require.False(t, done)
		d, done = it.Next()
		require.True(t, time.Millisecond*2900 <= d && d <= time.Millisecond*3100)
		require.False(t, done)
		d, done = it.Next()
		require.Equal(t, time.Duration(0), d)
		require.True(t, done)
	}

	// for test coverage
	b = NewLinear(time.Millisecond * -100)
	b = WithJitter(time.Millisecond * 100)(b)
	it := b.Iterator()
	d, done := it.Next()
	require.Equal(t, time.Duration(0), d)
	require.False(t, done)
	d, done = it.Next()
	require.Equal(t, time.Duration(0), d)
	require.False(t, done)
}
