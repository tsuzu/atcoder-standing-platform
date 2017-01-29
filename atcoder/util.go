package atcoder

import "time"

func IsInvalidFinishTime(t time.Time) bool {
	return t.Equal(invalidFinishTime)
}
