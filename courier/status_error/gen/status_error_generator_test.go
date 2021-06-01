package gen

import (
	"os"
	"testing"

	"github.com/artisanhe/tools/codegen"
)

func init() {
	os.Chdir("..")
}

func TestGen(t *testing.T) {
	statusErrorGenerator := StatusErrorGenerator{}
	codegen.Generate(&statusErrorGenerator)
}
