package corellium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1TeamResource exercises the corellium_team resource's
// full CRUD lifecycle, including add/remove-user-from-team, against a mock
// API server. No TF_ACC, live token, or network egress beyond loopback.
func TestUnitCorelliumV1TeamResource(t *testing.T) {
	const teamID = "66666666-6666-6666-6666-666666666666"
	const userID = "77777777-7777-7777-7777-777777777777"

	var mu sync.Mutex
	created := false
	label := "test"
	userIDs := map[string]bool{}

	srv := newMockAPIServer(t)

	// Create checks for a label collision by listing teams first, so the
	// list handler must reflect current state, including before Create
	// has ever run.
	teamsList := func() map[string]any {
		mu.Lock()
		defer mu.Unlock()
		results := []map[string]any{}
		if created {
			users := make([]map[string]any, 0, len(userIDs))
			for id := range userIDs {
				users = append(users, map[string]any{"id": id, "label": "user"})
			}
			results = append(results, map[string]any{"id": teamID, "label": label, "users": users})
		}
		return map[string]any{
			"results": results,
			"count":   len(results),
			"total":   len(results),
			"limit":   50,
			"offset":  0,
			"hasMore": false,
		}
	}
	srv.handle("GET /api/v3/teams", func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(http.StatusOK, teamsList())(w, r)
	})
	srv.handle("POST /api/v3/teams", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		created = true
		mu.Unlock()
		jsonHandler(http.StatusOK, map[string]string{"id": teamID})(w, r)
	})
	srv.handle("PATCH /api/v3/teams/"+teamID, func(w http.ResponseWriter, r *http.Request) {
		var patch map[string]any
		_ = json.NewDecoder(r.Body).Decode(&patch)
		if name, ok := patch["name"].(string); ok {
			mu.Lock()
			label = name
			mu.Unlock()
		}
		jsonHandler(http.StatusOK, map[string]any{})(w, r)
	})
	srv.handle("PUT /api/v3/teams/"+teamID+"/users/"+userID, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		userIDs[userID] = true
		mu.Unlock()
		jsonHandler(http.StatusOK, map[string]any{})(w, r)
	})
	srv.handle("DELETE /api/v3/teams/"+teamID+"/users/"+userID, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		delete(userIDs, userID)
		mu.Unlock()
		jsonHandler(http.StatusOK, map[string]any{})(w, r)
	})
	srv.handle("DELETE /api/v3/teams/"+teamID, jsonHandler(http.StatusOK, map[string]any{}))

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_team" "test" {
                    label = "test"
                    users = []
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_team.test", "id"),
					resource.TestCheckResourceAttr("corellium_team.test", "label", "test"),
					resource.TestCheckResourceAttr("corellium_team.test", "users.#", "0"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
                resource "corellium_team" "test" {
                    label = "test_updated"
                    users = [{ id = %q }]
                }`, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_team.test", "label", "test_updated"),
					resource.TestCheckResourceAttr("corellium_team.test", "users.#", "1"),
					resource.TestCheckResourceAttr("corellium_team.test", "users.0.id", userID),
				),
			},
		},
	})
}
