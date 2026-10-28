package corellium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1UserResource exercises the corellium_user resource's
// full CRUD lifecycle against a mock API server, with no TF_ACC, live
// token, or network egress beyond loopback.
func TestUnitCorelliumV1UserResource(t *testing.T) {
	const userID = "11111111-1111-1111-1111-111111111111"

	// The mock user record is mutable: PATCH updates it, and the teams
	// list (the only way this resource's Read can find a user, since
	// there's no dedicated get-user-by-id endpoint) always reflects the
	// latest values.
	var mu sync.Mutex
	user := map[string]any{
		"id":            userID,
		"label":         "ACC TEST ONE",
		"name":          "ACCTESTONENAME",
		"email":         "ACCTESTEMAILONE@ai.moda",
		"administrator": false,
	}

	srv := newMockAPIServer(t)
	srv.handle("POST /api/v3/users", jsonHandler(http.StatusOK, map[string]string{"id": userID}))
	srv.handle("PATCH /api/v3/users/"+userID, func(w http.ResponseWriter, r *http.Request) {
		var patch map[string]any
		_ = json.NewDecoder(r.Body).Decode(&patch)

		mu.Lock()
		for _, k := range []string{"label", "name", "email", "administrator"} {
			if v, ok := patch[k]; ok {
				user[k] = v
			}
		}
		mu.Unlock()

		jsonHandler(http.StatusOK, map[string]any{})(w, r)
	})
	srv.handle("DELETE /api/v3/users/"+userID, jsonHandler(http.StatusOK, map[string]any{}))
	srv.handle("GET /api/v3/teams", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		body := map[string]any{
			"results": []map[string]any{
				{"id": "all-users", "label": "All Users", "users": []map[string]any{user}},
			},
			"count":   1,
			"total":   1,
			"limit":   50,
			"offset":  0,
			"hasMore": false,
		}
		mu.Unlock()

		jsonHandler(http.StatusOK, body)(w, r)
	})

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resourceConfig := func(label, name, email string, administrator bool) string {
		return fmt.Sprintf(`
			resource "corellium_user" "test" {
				administrator = %t
				label         = %q
				name          = %q
				email         = %q
				password      = "fdsahj29sd8dh@#$"
			}`, administrator, label, name, email)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("ACC TEST ONE", "ACCTESTONENAME", "ACCTESTEMAILONE@ai.moda", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_user.test", "id"),
					resource.TestCheckResourceAttr("corellium_user.test", "label", "ACC TEST ONE"),
					resource.TestCheckResourceAttr("corellium_user.test", "name", "ACCTESTONENAME"),
					resource.TestCheckResourceAttr("corellium_user.test", "email", "ACCTESTEMAILONE@ai.moda"),
					resource.TestCheckResourceAttr("corellium_user.test", "administrator", "false"),
				),
			},
			{
				Config: providerConfig + resourceConfig("ACC TEST TWO", "ACCTESTONENAME", "ACCTESTEMAILTWO@ai.moda", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_user.test", "label", "ACC TEST TWO"),
					resource.TestCheckResourceAttr("corellium_user.test", "email", "ACCTESTEMAILTWO@ai.moda"),
				),
			},
		},
	})
}
