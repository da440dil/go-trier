# go-trier

[![Build Status](https://travis-ci.com/da440dil/go-trier.svg?branch=master)](https://travis-ci.com/da440dil/go-trier)
[![Coverage Status](https://coveralls.io/repos/github/da440dil/go-trier/badge.svg?branch=master)](https://coveralls.io/github/da440dil/go-trier?branch=master)
[![GoDoc](https://godoc.org/github.com/da440dil/go-trier?status.svg)](https://godoc.org/github.com/da440dil/go-trier)
[![Go Report Card](https://goreportcard.com/badge/github.com/da440dil/go-trier)](https://goreportcard.com/report/github.com/da440dil/go-trier)

(Re-)execution for Go functions (within configurable limits) with cancellation.

## Basic usage

```go
i := 0
// Create retriable function
fn := func() (bool, error) {
	// Return true if number is even
	ok := (i % 2) == 0
	i++
	return ok, nil
}
for j := 0; j < 3; j++ {
	// Create context for cancelling tries prematurely
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	// Execute function
	err := trier.Try(
		fn,
		// Set maximum number of retries
		trier.WithRetryCount(1),
		// Set delay between retries
		trier.WithRetryDelay(time.Millisecond*100),
		// Set maximum time randomly added to delay between retries
		trier.WithRetryJitter(time.Millisecond*20),
		// Set context
		trier.WithContext(ctx),
	)
	if err != nil {
		if err == trier.ErrRetryCountExceeded {
			// Failure
		} else {
			// Error
		}
	} else {
		// Success
	}
}
```

## Example usage

- [example](https://github.com/da440dil/go-counter/blob/master/examples/counter-with-retry/main.go) usage with [rate limiting](https://github.com/da440dil/go-counter)
- [example](https://github.com/da440dil/go-locker/blob/master/examples/locker-with-retry/main.go) usage with [locking](https://github.com/da440dil/go-locker)