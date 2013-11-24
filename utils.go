package gox

import (
	"time"
)

func timeFromUnixMicro(micro int64) time.Time {
	nanos := micro * 1e3
	seconds := int64(nanos / 1e9)
	nanos = int64(nanos % 1e9)

	return time.Unix(seconds, nanos)
}
