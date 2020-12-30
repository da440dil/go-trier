package main

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
