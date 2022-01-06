# go-trier

[![Build Status](https://travis-ci.com/da440dil/go-trier.svg?branch=master)](https://travis-ci.com/da440dil/go-trier)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-trier/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-trier?branch=master)
[![Go Reference](https://pkg.go.dev/badge/github.com/da440dil/go-trier.svg)](https://pkg.go.dev/github.com/da440dil/go-trier)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-trier)](https://goreportcard.com/report/github.com/da440dil/go-trier)

Re-execution for functions within configurable limits.

[Example](./examples/jitter/main.go) usage:

```go
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/da440dil/go-trier"
)

func main() {
	// Create trier.
	tr := trier.NewTrier(
		// Use linear algorithm to create delay between retries.
		trier.Linear(time.Millisecond*10),
		// Set maximum number of retries.
		trier.WithMaxRetries(5),
		// Set maximum duration randomly added to or extracted from delay
		// between retries to improve performance under high contention
		trier.WithJitter(time.Millisecond*5),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()
	for i := 0; i < 3; i++ {
		ok, err := tr.Try(ctx, func(ctx context.Context) (bool, error) {
			return i == 0, nil
		})
		fmt.Printf("{ i: %v, ok: %v, err: %v }\n", i, ok, err)
	}
	// Output:
	// { i: 0, ok: true, err: <nil> }
	// { i: 1, ok: false, err: <nil> }
	// { i: 2, ok: false, err: context deadline exceeded }
}
```
