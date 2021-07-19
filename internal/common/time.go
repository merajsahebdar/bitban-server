package common

import (
	"time"

	"github.com/volatiletech/null/v8"
)

// MayCurrentUtc
func MayCurrentUtc() null.Time {
	return null.TimeFrom(time.Now().UTC())
}
