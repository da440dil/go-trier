package trier_test

import (
	"context"
	"fmt"
	"time"

	"github.com/da440dil/go-trier"
)

func ExampleTry() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Execute function
		if err := trier.Try(fn); err != nil {
			if _, ok := err.(trier.TTLError); !ok {
				fmt.Println(err)
			}
		}
	}
	// Output:
	// 0
	// 2
}

func ExampleTry_withRetryCount() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Execute function
		// Set maximum number of retries
		err := trier.Try(fn, trier.WithRetryCount(1))
		if err != nil {
			fmt.Println(err)
		}
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleTry_withRetryDelay() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Execute function
		err := trier.Try(
			fn,
			// Set maximum number of retries
			trier.WithRetryCount(1),
			// Set delay between retries
			trier.WithRetryDelay(time.Millisecond*10),
		)
		if err != nil {
			fmt.Println(err)
		}
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleTry_withRetryJitter() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Execute function
		err := trier.Try(
			fn,
			// Set maximum number of retries
			trier.WithRetryCount(1),
			// Set delay between retries
			trier.WithRetryDelay(time.Millisecond*10),
			// Set maximum time randomly added to delay between retries
			trier.WithRetryJitter(time.Millisecond*5),
		)
		if err != nil {
			fmt.Println(err)
		}
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleTrier_Try() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Create new trier
		t, err := trier.New()
		if err != nil {
			fmt.Println(err)
		}
		// Execute function
		if err = t.Try(fn); err != nil {
			if _, ok := err.(trier.TTLError); !ok {
				fmt.Println(err)
			}
		}
	}
	// Output:
	// 0
	// 2
}

func ExampleTrier_Try_withRetryCount() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Create new trier
		// Set maximum number of retries
		t, err := trier.New(trier.WithRetryCount(1))
		if err != nil {
			fmt.Println(err)
		}
		// Execute function
		if err = t.Try(fn); err != nil {
			fmt.Println(err)
		}
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleTrier_Try_withRetryDelay() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Create new trier
		t, err := trier.New(
			// Set maximum number of retries
			trier.WithRetryCount(1),
			// Set delay between retries
			trier.WithRetryDelay(time.Millisecond*10),
		)
		if err != nil {
			fmt.Println(err)
		}
		// Execute function
		if err = t.Try(fn); err != nil {
			fmt.Println(err)
		}
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleTrier_Try_withRetryJitter() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		// Create new trier
		t, err := trier.New(
			// Set maximum number of retries
			trier.WithRetryCount(1),
			// Set delay between retries
			trier.WithRetryDelay(time.Millisecond*10),
			// Set maximum time randomly added to delay between retries
			trier.WithRetryJitter(time.Millisecond*5),
		)
		if err != nil {
			fmt.Println(err)
		}
		// Execute function
		if err = t.Try(fn); err != nil {
			fmt.Println(err)
		}
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleTrier_Try_withContext() {
	i := 0
	// Create retriable function
	fn := func() (bool, time.Duration, error) {
		ok := (i % 2) == 0
		if ok {
			fmt.Println(i)
		}
		i++
		return ok, -1, nil
	}
	for j := 0; j < 3; j++ {
		var timeout time.Duration
		if (j % 2) == 0 {
			timeout = time.Millisecond * 200
		} else {
			timeout = time.Millisecond * 20
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Create new trier
		t, err := trier.New(
			// Set maximum number of retries
			trier.WithRetryCount(1),
			// Set delay between retries
			trier.WithRetryDelay(time.Millisecond*100),
			// Set context for cancelling tries prematurely
			trier.WithContext(ctx),
		)
		if err != nil {
			fmt.Println(err)
		}

		// Execute function
		if err = t.Try(fn); err != nil {
			if err == context.DeadlineExceeded {
				fmt.Println("X")
			} else {
				fmt.Println(err)
			}

		}
	}
	// Output:
	// 0
	// X
	// 2
}
