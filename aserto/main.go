package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/aserto-dev/go-aserto"
	"github.com/aserto-dev/go-aserto/az"
	"github.com/aserto-dev/go-aserto/middleware"
	"github.com/aserto-dev/go-aserto/middleware/gorillaz"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func AsertoAuthorizer(addr, tenantID, apiKey, policy string) (*gorillaz.Middleware, error) {
	azClient, err := az.New(
		aserto.WithAddr(addr),
		aserto.WithTenantID(tenantID),
		aserto.WithAPIKeyAuth(apiKey),
	)
	if err != nil {
		return nil, err
	}

	mw := gorillaz.New(
		azClient,
		&middleware.Policy{
			Name:     policy,
			Decision: "allowed",
		},
	).WithPolicyFromURL("gorbac")
	mw.Identity.Mapper(func(r *http.Request, identity middleware.Identity) {
		if username, _, ok := r.BasicAuth(); ok {
			identity.Subject().ID(username)
		}
	})
	return mw, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	authorizerAddr := os.Getenv("AUTHORIZER_ADDRESS")
	if authorizerAddr == "" {
		authorizerAddr = "authorizer.prod.aserto.com:8443"
	}
	apiKey := os.Getenv("AUTHORIZER_API_KEY")
	policy := os.Getenv("POLICY_NAME")
	tenantID := os.Getenv("TENANT_ID")

	authorizer, err := AsertoAuthorizer(authorizerAddr, tenantID, apiKey, policy)
	if err != nil {
		log.Fatal("Failed to create authorizer:", err)
	}

	log.Print(os.Getenv("AUTHORIZER_API_KEY"))

	router := mux.NewRouter()
	router.HandleFunc("/api/{asset}", server.Handler).Methods("GET", "POST", "DELETE")
	router.Use(authorizer.Handler)

	server.Start(router)
}
