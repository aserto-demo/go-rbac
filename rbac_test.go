package rbac_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aserto-demo/go-rbac/pkg/server"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type test struct {
	user     string
	action   string
	resource string
	expected bool
}

func (t *test) name() string {
	return fmt.Sprintf("%s:%s:%s", t.user, t.action, t.resource)
}

func (t *test) url() string {
	return fmt.Sprintf("http://localhost:%d/api/%s", server.Port, t.resource)
}

func (t *test) run(tt *testing.T) {
	req, err := http.NewRequest(t.action, t.url(), http.NoBody)
	require.NoError(tt, err, "failed to create request for: %v", t)

	req.SetBasicAuth(t.user, "x")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(tt, err, "failed to make http request: %v", req)
	defer resp.Body.Close()

	expected := lo.Ternary(t.expected, http.StatusOK, http.StatusForbidden)
	assert.Equal(tt, expected, resp.StatusCode)

}

var tests = []test{
	{"summer@the-smiths.com", "GET", "mega-seed", false},
	{"summer@the-smiths.com", "GET", "portal-gun", false},
	{"summer@the-smiths.com", "GET", "space-cruiser", true},
	{"summer@the-smiths.com", "PUT", "mega-seed", false},
	{"summer@the-smiths.com", "PUT", "portal-gun", false},
	{"summer@the-smiths.com", "PUT", "space-cruiser", false},
	{"summer@the-smiths.com", "DELETE", "mega-seed", false},
	{"summer@the-smiths.com", "DELETE", "portal-gun", false},
	{"summer@the-smiths.com", "DELETE", "space-cruiser", false},

	{"morty@the-citadel.com", "GET", "mega-seed", true},
	{"morty@the-citadel.com", "GET", "portal-gun", true},
	{"morty@the-citadel.com", "GET", "space-cruiser", true},
	{"morty@the-citadel.com", "PUT", "mega-seed", true},
	{"morty@the-citadel.com", "PUT", "portal-gun", true},
	{"morty@the-citadel.com", "PUT", "space-cruiser", true},
	{"morty@the-citadel.com", "DELETE", "mega-seed", true},
	{"morty@the-citadel.com", "DELETE", "portal-gun", true},
	{"morty@the-citadel.com", "DELETE", "space-cruiser", false},

	{"rick@the-citadel.com", "GET", "mega-seed", true},
	{"rick@the-citadel.com", "GET", "portal-gun", true},
	{"rick@the-citadel.com", "GET", "space-cruiser", true},
	{"rick@the-citadel.com", "PUT", "mega-seed", false},
	{"rick@the-citadel.com", "PUT", "portal-gun", true},
	{"rick@the-citadel.com", "PUT", "space-cruiser", true},
	{"rick@the-citadel.com", "DELETE", "mega-seed", false},
	{"rick@the-citadel.com", "DELETE", "portal-gun", false},
	{"rick@the-citadel.com", "DELETE", "space-cruiser", true},
}

func TestRBAC(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name(), test.run)
	}
}
