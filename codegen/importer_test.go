package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPkgImportPathAndExpose(t *testing.T) {
	cases := []struct {
		full       string
		importPath string
		expose     string
	}{
		{
			"github.com/stretchr/testify/assert",
			"github.com/stretchr/testify/assert",
			"",
		},
		{
			"golib/tes.Test",
			"golib/tes",
			"Test",
		},
		{
			"github.com/stretchr/testify/assert.Equal",
			"github.com/stretchr/testify/assert",
			"Equal",
		},
		{
			"git.chinawayltd.com/stretchr/testify/assert.Equal",
			"git.chinawayltd.com/stretchr/testify/assert",
			"Equal",
		},
		{
			"gopkg.in/yaml.v2.Marshal",
			"gopkg.in/yaml.v2",
			"Marshal",
		},
	}

	for _, c := range cases {
		importPath, expose := getPkgImportPathAndExpose(c.full)
		assert.Equal(t, c.importPath, importPath)
		assert.Equal(t, c.expose, expose)
	}
}

func TestImports(t *testing.T) {
	tt := assert.New(t)
	imports := Importer{}

	tt.Equal(`[]string{"1", "2"}`, imports.Sdump([]string{"1", "2"}))
	tt.Equal(`[]interface {}{"1", nil}`, imports.Sdump([]interface{}{"1", nil}))
	tt.Equal(`map[string]int{"1": 2}`, imports.Sdump(map[string]int{"1": 2}))
}
