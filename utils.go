package mtgox

import (
	"fmt"
	"strconv"
	"time"
)

func timeFromUnixMicro(micro int64) time.Time {
	nanos := micro * 1e3
	seconds := int64(nanos / 1e9)
	nanos = int64(nanos % 1e9)

	return time.Unix(seconds, nanos)
}

type EpochTime struct {
	time.Time
}

func (t *EpochTime) UnmarshalJSON(b []byte) error {
	nowInt, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		fmt.Printf("Unmarshal error: %s\n", err.Error())
		return err
	}

	t.Time = timeFromUnixMicro(nowInt)

	return nil
}

func formatVolume(volume int64) string {
	return fmt.Sprintf("%v", float64(volume)/BitcoinDivision)
}
