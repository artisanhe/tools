package gen

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/artisanhe/tools/courier/swagger/gen"

	"github.com/go-courier/oas"
	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"

	"github.com/artisanhe/tools/codegen"
	"github.com/artisanhe/tools/courier/client/gen/enums"
	"github.com/artisanhe/tools/courier/client/gen/v2"
	"github.com/artisanhe/tools/courier/client/gen/v3"
	"github.com/artisanhe/tools/courier/status_error"
)

type ClientGenerator struct {
	File             string
	SpecURL          string
	BaseClient       string
	ServiceName      string
	swagger          *spec.Swagger
	openAPI          *oas.OpenAPI
	statusErrCodeMap status_error.StatusErrorCodeMap
}

func (g *ClientGenerator) Load(cwd string) {
	if g.SpecURL == "" && g.File == "" {
		logrus.Panicf("missing spec-url or file")
		return
	}

	if g.SpecURL != "" {
		g.loadBySpecURL()
	}

	if g.File != "" {
		g.loadByFile()
	}
}

func (g *ClientGenerator) loadByFile() {
	data, err := ioutil.ReadFile(g.File)
	if err != nil {
		panic(err)
	}

	g.swagger, g.openAPI, g.statusErrCodeMap = bytesToSwaggerOrOpenAPI(data)
}

func (g *ClientGenerator) loadBySpecURL() {
	hc := http.Client{}
	req, err := http.NewRequest("GET", g.SpecURL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	g.swagger, g.openAPI, g.statusErrCodeMap = bytesToSwaggerOrOpenAPI(bodyBytes)
	if g.ServiceName == "" {
		g.ServiceName = getIdFromUrl(g.SpecURL)
	}
}

func bytesToSwaggerOrOpenAPI(data []byte) (*spec.Swagger, *oas.OpenAPI, status_error.StatusErrorCodeMap) {
	checker := struct {
		OpenAPI string `json:"openapi"`
	}{}

	err := json.Unmarshal(data, &checker)
	if err != nil {
		panic(err)
	}

	data = bytes.Replace(data, []byte("golib/timelib"), []byte("golib/tools/timelib"), -1)
	data = bytes.Replace(data, []byte("golib/httplib"), []byte("golib/tools/httplib"), -1)

	data = bytes.Replace(data, []byte("golib/tools"), []byte("github.com/artisanhe/tools"), -1)

	data = bytes.Replace(data, []byte("git.chinawayltd.com/git.chinawayltd.com"), []byte("git.chinawayltd.com"), -1)

	data = bytes.Replace(data, []byte("«"), []byte{}, -1)
	data = bytes.Replace(data, []byte("»"), []byte{}, -1)

	statusErrCodeMap := status_error.StatusErrorCodeMap{}

	regexp.MustCompile("@httpError[^;]+;").ReplaceAllFunc(data, func(i []byte) []byte {
		v := bytes.Replace(i, []byte(`\"`), []byte(`"`), -1)
		s := status_error.ParseString(string(v))
		statusErrCodeMap[s.Code] = *s
		return i
	})

	if checker.OpenAPI == "" {
		swagger := new(spec.Swagger)

		err := json.Unmarshal(data, swagger)
		if err != nil {
			panic(err)
		}

		return swagger, nil, statusErrCodeMap
	}

	openAPI := new(oas.OpenAPI)

	err = json.Unmarshal(data, openAPI)
	if err != nil {
		panic(err)
	}

	return nil, openAPI, statusErrCodeMap
}

func getIdFromUrl(url string) string {
	paths := strings.Split(url, "/")
	return paths[len(paths)-1]
}

func (g *ClientGenerator) compatibleWithOuterSystem() {
	if g.swagger == nil {
		return
	}
	for _, pathItem := range g.swagger.Paths.Paths {
		if pathItem.PathItemProps.Get != nil {
			pathItem.PathItemProps.Get.ID = smallCamelToBigCamel(pathItem.PathItemProps.Get.ID)
		}
		if pathItem.PathItemProps.Post != nil {
			pathItem.PathItemProps.Post.ID = smallCamelToBigCamel(pathItem.PathItemProps.Post.ID)
		}
		if pathItem.PathItemProps.Put != nil {
			pathItem.PathItemProps.Put.ID = smallCamelToBigCamel(pathItem.PathItemProps.Put.ID)
		}
		if pathItem.PathItemProps.Delete != nil {
			pathItem.PathItemProps.Delete.ID = smallCamelToBigCamel(pathItem.PathItemProps.Delete.ID)
		}
		if pathItem.PathItemProps.Head != nil {
			pathItem.PathItemProps.Head.ID = smallCamelToBigCamel(pathItem.PathItemProps.Head.ID)
		}
		if pathItem.PathItemProps.Patch != nil {
			pathItem.PathItemProps.Patch.ID = smallCamelToBigCamel(pathItem.PathItemProps.Patch.ID)
		}
		if pathItem.PathItemProps.Options != nil {
			pathItem.PathItemProps.Options.ID = smallCamelToBigCamel(pathItem.PathItemProps.Options.ID)
		}
	}

	for _, def := range g.swagger.Definitions {
		for propKey, prop := range def.SchemaProps.Properties {
			if prop.Extensions == nil {
				prop.Extensions = make(map[string]interface{})
			}
			if _, ok := prop.Extensions[gen.XField]; !ok {
				prop.Extensions[gen.XField] = smallCamelToBigCamel(propKey)
			}
			def.SchemaProps.Properties[propKey] = prop
		}
	}
}

func smallCamelToBigCamel(old string) string {
	return strings.ToUpper(string(old[0])) + old[1:]
}

func (g *ClientGenerator) Pick() {
	g.compatibleWithOuterSystem()
}

func (g *ClientGenerator) Output(cwd string) codegen.Outputs {
	outputs := codegen.Outputs{}
	pkgName := codegen.ToLowerSnakeCase("Client-" + g.ServiceName)

	if g.swagger != nil {
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "client.go")), v2.ToClient(g.BaseClient, g.ServiceName, g.swagger))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "types.go")), v2.ToTypes(g.ServiceName, pkgName, g.swagger))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "enums.go")), enums.ToEnums(g.ServiceName, pkgName))
	}

	if g.openAPI != nil {
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "client.go")), v3.ToClient(g.BaseClient, g.ServiceName, g.openAPI))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "types.go")), v3.ToTypes(g.ServiceName, pkgName, g.openAPI))
		outputs.Add(codegen.GeneratedSuffix(path.Join(pkgName, "enums.go")), enums.ToEnums(g.ServiceName, pkgName))
	}

	return outputs
}
