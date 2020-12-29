package trier

import (
	"context"
	"math"
	"testing"
	"time"
)

func Benchmark_Iterator(b *testing.B) {
	d := time.Second
	tests := map[string]Iterator{
		"Constant":        NewConstant(d).Iterator(),
		"Linear":          NewLinear(d).Iterator(),
		"LinearRate":      NewLinearRate(d, d).Iterator(),
		"Exponential":     NewExponential(d).Iterator(),
		"ExponentialRate": NewExponentialRate(d, 1).Iterator(),
	}
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Next()
			}
		})
	}
}

func Benchmark_IteratorWithMaxRetries(b *testing.B) {
	d := time.Second
	fn := WithMaxRetries(math.MaxInt64)
	tests := map[string]Iterator{
		"Constant":        fn(NewConstant(d)).Iterator(),
		"Linear":          fn(NewLinear(d)).Iterator(),
		"LinearRate":      fn(NewLinearRate(d, d)).Iterator(),
		"Exponential":     fn(NewExponential(d)).Iterator(),
		"ExponentialRate": fn(NewExponentialRate(d, 1)).Iterator(),
	}
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Next()
			}
		})
	}
}

func Benchmark_IteratorWithJitter(b *testing.B) {
	d := time.Second
	fn := WithJitter(time.Millisecond * 100)
	tests := map[string]Iterator{
		"Constant":        fn(NewConstant(d)).Iterator(),
		"Linear":          fn(NewLinear(d)).Iterator(),
		"LinearRate":      fn(NewLinearRate(d, d)).Iterator(),
		"Exponential":     fn(NewExponential(d)).Iterator(),
		"ExponentialRate": fn(NewExponentialRate(d, 1)).Iterator(),
	}
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Next()
			}
		})
	}
}

func Benchmark_Iterable(b *testing.B) {
	d := time.Second
	tests := map[string]Iterable{
		"Constant":        NewConstant(d),
		"Linear":          NewLinear(d),
		"LinearRate":      NewLinearRate(d, d),
		"Exponential":     NewExponential(d),
		"ExponentialRate": NewExponentialRate(d, 1),
	}
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Iterator()
			}
		})
	}
}

func Benchmark_IterableWithMaxRetries(b *testing.B) {
	d := time.Second
	fn := WithMaxRetries(math.MaxInt64)
	tests := map[string]Iterable{
		"Constant":        fn(NewConstant(d)),
		"Linear":          fn(NewLinear(d)),
		"LinearRate":      fn(NewLinearRate(d, d)),
		"Exponential":     fn(NewExponential(d)),
		"ExponentialRate": fn(NewExponentialRate(d, 1)),
	}
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Iterator()
			}
		})
	}
}

func Benchmark_IterableWithJitter(b *testing.B) {
	d := time.Second
	fn := WithJitter(time.Millisecond * 100)
	tests := map[string]Iterable{
		"Constant":        fn(NewConstant(d)),
		"Linear":          fn(NewLinear(d)),
		"LinearRate":      fn(NewLinearRate(d, d)),
		"Exponential":     fn(NewExponential(d)),
		"ExponentialRate": fn(NewExponentialRate(d, 1)),
	}
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Iterator()
			}
		})
	}
}

type mock bool

func (m *mock) Try(ctx context.Context) (bool, error) {
	v := *m
	*m = !v
	return bool(v), nil
}

func Benchmark_Trier(b *testing.B) {
	d := time.Millisecond * 10
	tests := map[string]Trier{
		"Constant":        NewTrier(NewConstant(d)),
		"Linear":          NewTrier(NewLinear(d)),
		"LinearRate":      NewTrier(NewLinearRate(d, d)),
		"Exponential":     NewTrier(NewExponential(d)),
		"ExponentialRate": NewTrier(NewExponentialRate(d, 1)),
	}
	ctx := context.Background()
	m := new(mock)
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Try(ctx, m.Try)
			}
		})
	}
}

func tryf(ctx context.Context) (bool, error) {
	return false, nil
}

func Benchmark_TrierWithMaxRetries(b *testing.B) {
	d := time.Millisecond * 10
	fn := WithMaxRetries(1)
	tests := map[string]Trier{
		"Constant":        NewTrier(NewConstant(d), fn),
		"Linear":          NewTrier(NewLinear(d), fn),
		"LinearRate":      NewTrier(NewLinearRate(d, d), fn),
		"Exponential":     NewTrier(NewExponential(d), fn),
		"ExponentialRate": NewTrier(NewExponentialRate(d, 1), fn),
	}
	ctx := context.Background()
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Try(ctx, tryf)
			}
		})
	}
}

func Benchmark_TrierWithJitter(b *testing.B) {
	d := time.Millisecond * 10
	fn := WithJitter(time.Millisecond * 2)
	tests := map[string]Trier{
		"Constant":        NewTrier(NewConstant(d), fn),
		"Linear":          NewTrier(NewLinear(d), fn),
		"LinearRate":      NewTrier(NewLinearRate(d, d), fn),
		"Exponential":     NewTrier(NewExponential(d), fn),
		"ExponentialRate": NewTrier(NewExponentialRate(d, 1), fn),
	}
	ctx := context.Background()
	m := new(mock)
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Try(ctx, m.Try)
			}
		})
	}
}

func Benchmark_TrierWithMaxRetriesJitter(b *testing.B) {
	d := time.Millisecond * 10
	fn1 := WithMaxRetries(1)
	fn2 := WithJitter(time.Millisecond * 2)
	tests := map[string]Trier{
		"Constant":        NewTrier(NewConstant(d), fn1, fn2),
		"Linear":          NewTrier(NewLinear(d), fn1, fn2),
		"LinearRate":      NewTrier(NewLinearRate(d, d), fn1, fn2),
		"Exponential":     NewTrier(NewExponential(d), fn1, fn2),
		"ExponentialRate": NewTrier(NewExponentialRate(d, 1), fn1, fn2),
	}
	ctx := context.Background()
	for name, tc := range tests {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tc.Try(ctx, tryf)
			}
		})
	}
}
