package main

import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-trier"
)

func main() {
	// Create trier.
	tr := trier.NewTrier(
		// Use exponential algorithm to create delay between retries.
		trier.Exponential(time.Millisecond*10),
		// Set maximum number of retries.
		trier.WithMaxRetries(5),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*400)
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
