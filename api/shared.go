package api

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
)

const Title = "Huma demo service"

type VersionConfig struct {
	Tag    string
	SemVer string
}

func (vc VersionConfig) Prefixer(op huma.Operation) huma.Operation {
	op.Path = fmt.Sprintf("/api/%s%s", vc.Tag, op.Path)
	return op
}

type Impl struct{}

func Setup[Impl any](router *gin.Engine, vc VersionConfig) {

	prefix := fmt.Sprintf("/api/%s", vc.Tag)
	config := huma.DefaultConfig(Title, vc.SemVer)
	config.DocsPath = fmt.Sprintf("%s/docs", prefix)
	config.OpenAPIPath = fmt.Sprintf("%s/docs/openapi", prefix)
	config.SchemasPath = fmt.Sprintf("%s/schemas", prefix)
	api := humagin.New(router, config)
	var implPtr *Impl
	implType := reflect.TypeOf(implPtr)
	args := []reflect.Value{
		reflect.ValueOf(implPtr),
		reflect.ValueOf(api),
		reflect.ValueOf(vc),
	}
	for i := 0; i < implType.NumMethod(); i++ {
		m := implType.Method(i)
		if strings.HasPrefix(m.Name, "Register") && len(m.Name) > 8 {
			m.Func.Call(args)
		}
	}
}
