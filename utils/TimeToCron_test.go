package utils

import (
	"testing"
	"time"
)

func TestTimeToCron(t *testing.T) {
	cronTime := TimeToCron(time.Date(2021, time.April, 24, 3, 11, 12, 0, time.Now().Location()))

	if cronTime != "12 11 3 24 4 *" {
		t.Errorf("Cron time is incorrect. Currently: %s", cronTime)
	}
}
