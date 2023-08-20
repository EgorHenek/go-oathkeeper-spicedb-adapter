package domain

import (
	"net/http"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/go-chi/render"
)

type CheckPermissionRequest v1.CheckPermissionRequest

type CheckPermissionResponse struct {
	Allowed bool `json:"allowed"`
}

func (rd *CheckPermissionResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	if rd.Allowed {
		render.Status(r, http.StatusOK)
	} else {
		render.Status(r, http.StatusForbidden)
	}
	return nil
}

func NewCheckPermissionResponse(allowed bool) render.Renderer {
	return &CheckPermissionResponse{
		Allowed: allowed,
	}
}

func (request *CheckPermissionRequest) Bind(_ *http.Request) error {
	return nil
}
