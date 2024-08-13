package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aserto-demo/go-rbac/pkg/authz"
	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-dev/go-aserto"
	"github.com/aserto-dev/go-aserto/az"
	"github.com/aserto-dev/go-aserto/middleware"
	"github.com/aserto-dev/go-aserto/middleware/gorillaz"
	"github.com/gorilla/mux"
)

func AsertoAuthorizer(addr string) (*gorillaz.Middleware, error) {
	azClient, err := az.New(aserto.WithAddr(addr))
	if err != nil {
		return nil, err
	}

	mw := gorillaz.New(
		azClient,
		&middleware.Policy{
			Decision: "allowed",
			Root:     "rebac",
		},
	)
	mw.Identity.Mapper(func(r *http.Request, identity middleware.Identity) {
		if username, _, ok := r.BasicAuth(); ok {
			identity.Subject().ID(username)
		}
	})
	return mw, nil
}

func main() {
	authorizerAddr := os.Getenv("AUTHORIZER_ADDRESS")
	if authorizerAddr == "" {
		authorizerAddr = "localhost:8282" // default topaz authorizer port
	}

	authorizer, err := AsertoAuthorizer(authorizerAddr)
	if err != nil {
		log.Fatal("Failed to create authorizer:", err)
	}

	log.Print(os.Getenv("AUTHORIZER_API_KEY"))

	router := mux.NewRouter()
	router.Use(authorizer.Check(
		gorillaz.WithObjectType("resource"),
		gorillaz.WithObjectIDFromVar("resource"),
		gorillaz.WithRelationMapper(authz.ActionFromMethod),
	).Handler)

	router.HandleFunc("/api/{resource}", server.Handler).Methods("GET", "PUT", "DELETE")

	server.Start(router)
}
