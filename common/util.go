package common

import (
	"time"

	"github.com/go-shadow/moment"
)

// GetCurrentDateStr time.Now to string
func GetCurrentDateStr() string {
	return moment.New().Format(TimeFormat)
}

// GetCurrentTimestamp take time zone into consideration
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func GetSecondsDiff(t1 time.Time, t2 time.Time) int64 {
	return t1.Unix() - t2.Unix()
}
