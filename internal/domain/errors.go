package domain

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	AppCode    int64  `json:"code,omitempty"`
	StatusText string `json:"status"`
	ErrorText  string `json:"error,omitempty"`
}

func (e ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,

		StatusText: "Bad Request",
		ErrorText:  err.Error(),
	}
}
