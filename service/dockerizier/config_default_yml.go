package dockerizier

import (
	"gopkg.in/yaml.v2"

	"github.com/artisanhe/tools/conf"
)

func toConfigDefaultYML(envVars conf.EnvVars) string {
	e := make(map[string]string)

	e["GOENV"] = "DEV"

	for key, envVar := range envVars {
		// the CanConfig is useless, now will expose all config values
		if envVar.CanConfig {
			if envVar.FallbackValue != nil {
				e[key] = envVar.GetFallbackValue(false)
			} else {
				e[key] = envVar.GetValue(false)
			}
		}
	}

	bytes, _ := yaml.Marshal(e)
	return string(bytes)
}
