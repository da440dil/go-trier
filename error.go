package trier

import "time"

type trierError string

func (e trierError) Error() string {
	return string(e)
}

// ErrInvalidRetryCount is the error returned when WithCounter receives invalid value.
const ErrInvalidRetryCount = trierError("trier: number of retries must be greater than or equal to 0")

// ErrInvalidRetryDelay is the error returned when WithRetryDelay receives invalid value.
const ErrInvalidRetryDelay = trierError(
	"trier: delay between retries must be greater than or equal to 1 millisecond" +
		" and must be greater than or equal to jitter",
)

// ErrInvalidRetryJitter is the error returned when WithRetryJitter receives invalid value.
const ErrInvalidRetryJitter = trierError(
	"trier: retry jitter must be greater than or equal to 1 millisecond" +
		" and must be less than or equal to delay",
)

// ErrTooManyRetries is the error wrapped with TTLError.
const ErrTooManyRetries = trierError("trier: number of retries exceeded")

// TTLError is the error returned when Counter failed to count.
type TTLError struct {
	err error
	ttl time.Duration
}

func newTTLError(ttl time.Duration) *TTLError {
	return &TTLError{
		err: ErrTooManyRetries,
		ttl: ttl,
	}
}

func (e *TTLError) Error() string {
	return e.err.Error()
}

// TTL returns TTL of a key.
func (e *TTLError) TTL() time.Duration {
	return e.ttl
}

func (e *TTLError) Unwrap() error {
	return e.err
}
