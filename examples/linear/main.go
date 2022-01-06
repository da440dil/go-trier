package main

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
		trier.Linear(time.Millisecond*100),
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
