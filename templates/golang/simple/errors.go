package main

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	MessageText string `json:"message"` // user-level status message
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		MessageText:    err.Error(),
	}
}

func ErrConflictRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		MessageText:    err.Error(),
	}
}

func ErrBadRequest(errMsg string) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: http.StatusBadRequest,
		MessageText:    errMsg,
	}
}

func ErrNotFoundStr(errMsg string) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: http.StatusNotFound,
		MessageText:    errMsg,
	}
}

var (
	ErrInvalidToken     = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, MessageText: "Invalid Token."}
	ErrIdentityNotFound = &ErrResponse{HTTPStatusCode: http.StatusUnauthorized, MessageText: "Identity not found."}
	ErrNotFound         = &ErrResponse{HTTPStatusCode: http.StatusNotFound, MessageText: "Resource not found."}
	ErrUserNotAdmin     = &ErrResponse{HTTPStatusCode: http.StatusForbidden, MessageText: "User not tenant admin."}
	ErrUserNotEditor    = &ErrResponse{HTTPStatusCode: http.StatusForbidden, MessageText: "User not tenant admin or service-provider-admin."}
	ErrMethodNotAllowed = &ErrResponse{HTTPStatusCode: http.StatusServiceUnavailable, MessageText: "Method not allowed."}
	ErrTenantNotFound   = &ErrResponse{HTTPStatusCode: http.StatusForbidden, MessageText: "Tenant not found."}
	ErrMultipleFound    = &ErrResponse{HTTPStatusCode: http.StatusInternalServerError, MessageText: "Multiple tenants found with the same search string."}
)

func ErrInconsistentData(errMsg string) *ErrResponse {
	return &ErrResponse{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		MessageText:    errMsg,
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, ErrNotFound)
}
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, ErrMethodNotAllowed)
}
