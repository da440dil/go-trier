# go-trier

[![Build Status](https://travis-ci.com/da440dil/go-trier.svg?branch=master)](https://travis-ci.com/da440dil/go-trier)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-trier/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-trier?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/da440dil/go-trier.svg)](https://pkg.go.dev/github.com/da440dil/go-trier)
[![GoDoc](https://godoc.org/github.com/da440dil/go-trier?status.svg)](https://godoc.org/github.com/da440dil/go-trier)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-trier)](https://goreportcard.com/report/github.com/da440dil/go-trier)

Re-execution for functions within configurable limits.

[Basic](./examples/linear/main.go) usage:

```go
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/da440dil/go-trier"
)

func main() {
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
			defer wg.Done()
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
		}(i)
	}
	wg.Wait()
	// Output (may differ on each run because of jitter usage):
	// #2: success
	// #0: failure
	// #1: error: context deadline exceeded
}
```

[Example](./examples/fibonacci/main.go) usage with custom iterable:

```go
import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-trier"
)

func main() {
	tr := trier.NewTrier(
		// Use fibonacci growth algorithm to create delay between retries
		newFibonacci(time.Millisecond*10),
		// Set maximum number of retries
		trier.WithMaxRetries(5),
	)
	now := time.Now()
	ctx := context.Background()
	tr.Try(ctx, func(ctx context.Context) (bool, error) {
		return false, nil
	})
	fmt.Println(time.Since(now).Truncate(time.Millisecond * 10))
	// Output (10ms + 20ms + 30ms + 50ms + 80ms):
	// 190ms
}

type fibonacci struct {
	prev, curr time.Duration
}

func newFibonacci(d time.Duration) trier.Iterable {
	return fibonacci{curr: d}
}

func (i *fibonacci) Next() (time.Duration, bool) {
	i.prev, i.curr = i.curr, i.prev+i.curr
	return i.curr, false
}

func (i fibonacci) Iterator() trier.Iterator {
	return &fibonacci{curr: i.curr}
}
```

[Benchmarks](./benchmarks.md)
