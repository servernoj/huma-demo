package v1

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	. "github.com/servernoj/huma-demo/api"
)

type RawBody struct {
	Body []byte
}

type VersionedImpl Impl

func (impl *VersionedImpl) RegisterVersion(api huma.API, vc VersionConfig) {
	huma.Register(
		api,
		vc.Prefixer(
			huma.Operation{
				OperationID: "get-version",
				Summary:     "Returns version of the called API",
				Method:      http.MethodGet,
				Path:        "/version",
			},
		),
		func(ctx context.Context, input *struct{}) (*RawBody, error) {
			return &RawBody{
				Body: []byte(vc.SemVer),
			}, nil
		},
	)
}
