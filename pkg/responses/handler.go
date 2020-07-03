package responses

import (
	"net/http"

	l "github.com/bmaynard/apimock/pkg/utils/logger"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}

type MockResponseItem struct {
	Response   interface{} `json:"response"`
	StatusCode int
}

type Handler struct {
	MockResponses []*MockResponseItem
	H             func(e []*MockResponseItem, w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(h.MockResponses, w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			l.Log.Infof("HTTP %d - %s", e.Status(), e)
			http.Error(w, e.Error(), e.Status())
		default:
			// Any error types we don't specifically look out for default
			// to serving a HTTP 500
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}
