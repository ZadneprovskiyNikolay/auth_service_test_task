package httputils

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type httpError struct {
	Message string `json:"error,omitempty"`
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, r, http.StatusBadRequest, err)
}

func InternalError(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusInternalServerError, errors.New(""))
}

func UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, r, http.StatusUnauthorized, err)
}

func Error(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.Status(r, status)
	render.JSON(w, r, httpError{Message: err.Error()})
}
