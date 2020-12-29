# go-trier

[![Build Status](https://travis-ci.com/da440dil/go-trier.svg?branch=master)](https://travis-ci.com/da440dil/go-trier)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-trier/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-trier?branch=master)
[![GoDoc](https://godoc.org/github.com/da440dil/go-trier?status.svg)](https://godoc.org/github.com/da440dil/go-trier)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-trier)](https://goreportcard.com/report/github.com/da440dil/go-trier)

Re-execution for functions within configurable limits.

Basic usage:

```go
tr := trier.NewTrier(
	// Use linear growth algorithm to create delay between retries
	trier.NewLinear(time.Millisecond*100),
	// Set maximum number of retries
	trier.WithMaxRetries(1),
	// Set maximum duration randomly added to or extracted from delay
	// between retries to improve performance under high contention
	trier.WithJitter(time.Millisecond*50),
)
var wg sync.WaitGroup
j := 3
for i := 0; i < j; i++ {
	wg.Add(1)
	go func(i int) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()
		ok, err := tr.Try(ctx, func(i int) trier.Retriable {
			return func(ctx context.Context) (bool, error) {
				i++
				return i%j == 0, nil
			}
		}(i))
		if err != nil {
			fmt.Printf("#%v: error: %v\n", i, err)
		} else if ok {
			fmt.Printf("#%v: success\n", i)
		} else {
			fmt.Printf("#%v: failure\n", i)
		}
		wg.Done()
	}(i)
}
wg.Wait()
// Output:
// #2: success
// #0: failure
// #1: error: context deadline exceeded
```
