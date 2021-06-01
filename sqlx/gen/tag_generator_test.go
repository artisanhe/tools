package gen

import (
	"testing"

	"github.com/artisanhe/tools/codegen"
)

func TestTagGen(t *testing.T) {
	clientGenerator := TagGenerator{
		WithDefaults: true,
	}
	clientGenerator.StructNames = []string{"User", "User2"}
	codegen.Generate(&clientGenerator)
}
