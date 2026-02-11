package ntime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-lazyer/north/ntime"
)

func TestGetDayOfWeek(t *testing.T) {
	fmt.Print(ntime.GetDayOfWeek(time.Now(), 1, 1))
}
