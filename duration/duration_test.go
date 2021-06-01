package duration_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/artisanhe/tools/duration"
)

func TestNewCost(t *testing.T) {
	cost := duration.NewDuration()
	time.Sleep(1000 * time.Millisecond)
	fmt.Println(cost.Get())
}
