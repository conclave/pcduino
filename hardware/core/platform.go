package core

import (
	"time"
)

func Millis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func Micros() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func Delay(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}

func DelayMicrosends(us time.Duration) {
	time.Sleep(us * time.Microsecond)
}

func DelayShed(ms time.Duration) {
	time.Sleep(ms * time.Millisecond)
}

func DelayMicrosendsSched(us time.Duration) {
	time.Sleep(us * time.Microsecond)
}

// func PAbort(s string) { // perror(s) + abort();

// }
