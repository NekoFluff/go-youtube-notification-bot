package utils

import (
	"fmt"
	"time"
)

func TimeToCron(t time.Time) string {
	return fmt.Sprintf("%v %v %v %v %v *", t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()))
}
