package service_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

func TestClear_SuccessPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	obj := e.GET("/api/v1/user").
		Expect().
		Status(http.StatusOK).JSON().Object()

	_ = obj
}
