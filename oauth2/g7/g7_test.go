package g7

import (
	"reflect"
	"testing"

	"github.com/artisanhe/tools/conf"
)

func TestOAuth(t *testing.T) {
	o := OAuth{}
	o.Name = "test"
	rv := reflect.ValueOf(o)
	envVarsForDocker, _ := conf.CollectEnvVars(rv, "D")
	envVarsForDocker.Print()
}
