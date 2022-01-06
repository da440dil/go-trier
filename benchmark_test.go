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
		"Constant":        Constant(d).Iterator(),
		"Linear":          Linear(d).Iterator(),
		"LinearRate":      LinearRate(d, d).Iterator(),
		"Exponential":     Exponential(d).Iterator(),
		"ExponentialRate": ExponentialRate(d, 1).Iterator(),
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
		"Constant":        fn(Constant(d)).Iterator(),
		"Linear":          fn(Linear(d)).Iterator(),
		"LinearRate":      fn(LinearRate(d, d)).Iterator(),
		"Exponential":     fn(Exponential(d)).Iterator(),
		"ExponentialRate": fn(ExponentialRate(d, 1)).Iterator(),
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
		"Constant":        fn(Constant(d)).Iterator(),
		"Linear":          fn(Linear(d)).Iterator(),
		"LinearRate":      fn(LinearRate(d, d)).Iterator(),
		"Exponential":     fn(Exponential(d)).Iterator(),
		"ExponentialRate": fn(ExponentialRate(d, 1)).Iterator(),
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
		"Constant":        Constant(d),
		"Linear":          Linear(d),
		"LinearRate":      LinearRate(d, d),
		"Exponential":     Exponential(d),
		"ExponentialRate": ExponentialRate(d, 1),
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
		"Constant":        fn(Constant(d)),
		"Linear":          fn(Linear(d)),
		"LinearRate":      fn(LinearRate(d, d)),
		"Exponential":     fn(Exponential(d)),
		"ExponentialRate": fn(ExponentialRate(d, 1)),
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
		"Constant":        fn(Constant(d)),
		"Linear":          fn(Linear(d)),
		"LinearRate":      fn(LinearRate(d, d)),
		"Exponential":     fn(Exponential(d)),
		"ExponentialRate": fn(ExponentialRate(d, 1)),
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
		"Constant":        NewTrier(Constant(d)),
		"Linear":          NewTrier(Linear(d)),
		"LinearRate":      NewTrier(LinearRate(d, d)),
		"Exponential":     NewTrier(Exponential(d)),
		"ExponentialRate": NewTrier(ExponentialRate(d, 1)),
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
		"Constant":        NewTrier(Constant(d), fn),
		"Linear":          NewTrier(Linear(d), fn),
		"LinearRate":      NewTrier(LinearRate(d, d), fn),
		"Exponential":     NewTrier(Exponential(d), fn),
		"ExponentialRate": NewTrier(ExponentialRate(d, 1), fn),
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
		"Constant":        NewTrier(Constant(d), fn),
		"Linear":          NewTrier(Linear(d), fn),
		"LinearRate":      NewTrier(LinearRate(d, d), fn),
		"Exponential":     NewTrier(Exponential(d), fn),
		"ExponentialRate": NewTrier(ExponentialRate(d, 1), fn),
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
		"Constant":        NewTrier(Constant(d), fn1, fn2),
		"Linear":          NewTrier(Linear(d), fn1, fn2),
		"LinearRate":      NewTrier(LinearRate(d, d), fn1, fn2),
		"Exponential":     NewTrier(Exponential(d), fn1, fn2),
		"ExponentialRate": NewTrier(ExponentialRate(d, 1), fn1, fn2),
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
