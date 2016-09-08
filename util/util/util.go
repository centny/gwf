package util

import (
	"runtime"
	"time"
)

//Timestamp parsing time.Time to timestamp
func Timestamp(t time.Time) int64 {
	return t.Local().UnixNano() / 1e6
}

//Time parsing timestamp to time.Time
func Time(timestamp int64) time.Time {
	return time.Unix(0, timestamp*1e6)
}

//Now return current timestamp
func Now() int64 {
	return Timestamp(time.Now())
}

//NowSec return curent timestamp by seconds
func NowSec() int64 {
	return Timestamp(time.Now()) / 1000 * 1000
}

//CPU return CPU count sub 1, min is 1
func CPU() int {
	i := runtime.NumCPU()
	if i > 1 {
		i = i - 1
	}
	return i
}
