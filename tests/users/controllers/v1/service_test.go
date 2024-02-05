package service_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	baseUrl = "http://localhost:8080"
)

func TestClear_SuccessPath(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  baseUrl,
		Reporter: httpexpect.NewRequireReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	e.GET("/api/v1/user").
		Expect().
		Status(http.StatusUnauthorized)

	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	res := e.POST("/api/v1/auth/login").
		WithJSON(Login{"ivan@email.com", "123"}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	res.Path("$.data").Object().Keys().ContainsOnly("access_token")

	token := res.Path("$.data").Object().Value("access_token").String().Raw()

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	res = auth.GET("/api/v1/user").
		Expect().
		Status(http.StatusOK).JSON().Object()

	_ = res
}
