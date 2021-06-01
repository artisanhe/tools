package gen

import (
	"testing"

	"github.com/artisanhe/tools/codegen"
)

func TestGen(t *testing.T) {
	clientGenerator := SqlFuncGenerator{}
	clientGenerator.StructName = "User"
	clientGenerator.Database = "DBTest"
	clientGenerator.WithTableInterfaces = true
	codegen.Generate(&clientGenerator)
}
