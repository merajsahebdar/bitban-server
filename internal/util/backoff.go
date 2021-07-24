package util

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// NewExponentialBackoff Instantiates an exponential backoff with the provided.
// maximum elapsed time.
func NewExponentialBackoff(t time.Duration) *backoff.ExponentialBackOff {
	b := &backoff.ExponentialBackOff{
		InitialInterval:     backoff.DefaultInitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         backoff.DefaultMaxInterval,
		MaxElapsedTime:      t,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}

	b.Reset()

	return b
}
