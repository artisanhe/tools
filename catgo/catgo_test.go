package catgo

import (
	"math/rand"
	"testing"
	"time"

	"github.com/artisanhe/tools/catgo/cat-go/cat"
)

const (
	TestType = "Connection"
)

func Test_A(t *testing.T) {

	(&Catgo{
		Domain:      "xx.test",
		Debug:       true,
		HttpServers: []string{"test.dubhe.chinawayltd.com:2280"},
	}).Init()

	case1()
	time.Sleep(time.Second)

}

func case1() {

	for i := 0; i < 100; i++ {
		cat.LogEvent("XConnection", "event-000", cat.FAIL)
	}

	t := cat.NewTransaction("Connection", "mammmo")
	defer t.Complete()

	if rand.Int31n(100) == 0 {
		t.SetStatus(cat.FAIL)
	}

	t.AddData("foo", "bar")

	t.NewEvent(TestType, "event-1")
	t.Complete()

	if rand.Int31n(100) == 0 {
		t.LogEvent(TestType, "event-2", cat.FAIL)
	} else {
		t.LogEvent(TestType, "event-2")
	}
	t.LogEvent(TestType, "event-3", cat.SUCCESS, "k=v")

	t.SetDurationStart(time.Now().Add(-5 * time.Second))
	t.SetTime(time.Now().Add(-5 * time.Second))
	t.SetDuration(time.Millisecond * 500)

	time.Sleep(5 * time.Second)
	cat.Shutdown()
}
