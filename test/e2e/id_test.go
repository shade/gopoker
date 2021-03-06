package id_test

import (
	"gopoker/internal/identity"
	"gopoker/internal/identity/db"
	"gopoker/internal/server"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/gavv/httpexpect.v2"
)

func TestUserCreation(t *testing.T) {
	// Set up identity server
	address := ":8080"
	tmpfile, _ := ioutil.TempFile("", "test_db")
	db := db.NewFileDB(tmpfile.Name())
	id := identity.NewIDGen(db, "random")
	handler := server.Routes("", address, id, nil)

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// TEST user creation.
	e.POST("/users").
		WithForm(map[string]string{
			"username": "joe",
			"password": "",
		}).
		Expect().Status(http.StatusBadRequest)
	e.POST("/users").
		WithForm(map[string]string{
			"username": "",
			"password": "joe",
		}).
		Expect().Status(http.StatusBadRequest)
	e.POST("/users").
		WithForm(map[string]string{
			"username": "joe",
			"password": "joe",
		}).
		Expect().Status(http.StatusOK)

	// TEST token fetching
	req := e.POST("/users").
		WithForm(map[string]string{
			"username": "joe1",
			"password": "password",
		}).
		Expect().Status(http.StatusOK)
	// Don't allow overwrite
	e.POST("/users").
		WithForm(map[string]string{
			"username": "joe1",
			"password": "password",
		}).
		Expect().Status(http.StatusBadRequest)

	token := req.JSON().Object().Value("token").String().Raw()
	eAuth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	eAuth.GET("/tables").Expect().Status(http.StatusOK)
	e.GET("/tables").Expect().Status(http.StatusForbidden)
}
