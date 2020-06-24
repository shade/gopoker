package id_test

import (
	"fmt"
	"go_poker/internal/identity"
	"go_poker/internal/identity/db"
	"go_poker/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/gavv/httpexpect.v2"
)

func TestUserCreation(t *testing.T) {
	// Set up identity server
	address := ":8080"
	db := db.NewFileDB("./test.csv")
	id := identity.NewIDGen(db, "random")
	handler := server.Routes("", address, id)

	// run server using httptest
	server := httptest.NewServer(handler)
	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// TEST user creation.
	e.POST("/users").Expect().Status(http.StatusBadRequest)
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
			"username": "joe",
			"password": "password",
		}).
		Expect().Status(http.StatusOK)

	token := req.JSON().Object().Value("token").String().Raw()
	fmt.Printf("TOKEN %s\n", token)
	eAuth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	eAuth.GET("/tables").Expect().Status(http.StatusOK)
	e.GET("/tables").Expect().Status(http.StatusForbidden)
}