package unique_str_test

import (
	"testing"

	"github.com/artisanhe/tools/misc/unique_str"
)

func TestUniqueStr(t *testing.T) {
	unique_str := unique_str.GenerateUniqueIDToStr(uint64(1121314151617181910), 25)
	t.Logf("unique_str:%s", unique_str)
}

func TestRandIDStr(t *testing.T) {
	t.Logf(unique_str.New())
}
