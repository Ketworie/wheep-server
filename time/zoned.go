package time

import (
	"fmt"
	"time"
)

var Zoned = "2006-01-02T15:04:05.999Z[MST]"

type JSONTime struct {
	time.Time
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.UTC().Format(Zoned))
	return []byte(stamp), nil
}

func (t JSONTime) MarshalBSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.UTC().Format(Zoned))
	return []byte(stamp), nil
}
