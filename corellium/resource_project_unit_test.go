package corellium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1ProjectResource exercises the corellium_project
// resource's Create/Read/Update lifecycle (no users/teams/keys) against a
// mock API server. No TF_ACC, live token, or network egress beyond
// loopback.
func TestUnitCorelliumV1ProjectResource(t *testing.T) {
	const projectID = "33333333-3333-3333-3333-333333333333"

	var mu sync.Mutex
	created := false
	name := "test"
	version := float64(1)
	internetAccess := false
	dhcp := false
	cores := float64(1)

	projectBody := func() map[string]any {
		mu.Lock()
		defer mu.Unlock()
		return map[string]any{
			"id":   projectID,
			"name": name,
			"settings": map[string]any{
				"version":        version,
				"internetAccess": internetAccess,
				"dhcp":           dhcp,
			},
			"quotas": map[string]any{
				"cores":     cores,
				"instances": cores * 2.5,
				"ram":       cores * 6144,
			},
			"quotasUsed": map[string]any{},
		}
	}

	srv := newMockAPIServer(t)
	srv.handle("GET /api/v3/projects", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		results := []map[string]any{}
		mu.Unlock()
		if created {
			results = append(results, projectBody())
		}
		jsonHandler(http.StatusOK, map[string]any{
			"results": results, "count": len(results), "total": len(results), "limit": 50, "offset": 0, "hasMore": false,
		})(w, r)
	})
	srv.handle("POST /api/v3/projects", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		created = true
		if n, ok := body["name"].(string); ok {
			name = n
		}
		mu.Unlock()
		jsonHandler(http.StatusOK, projectBody())(w, r)
	})
	srv.handle("GET /api/v3/projects/"+projectID, func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(http.StatusOK, projectBody())(w, r)
	})
	srv.handle("PATCH /api/v3/projects/"+projectID, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		mu.Lock()
		if n, ok := body["name"].(string); ok {
			name = n
		}
		if s, ok := body["settings"].(map[string]any); ok {
			if v, ok := s["version"].(float64); ok {
				version = v
			}
			if v, ok := s["internetAccess"].(bool); ok {
				internetAccess = v
			}
			if v, ok := s["dhcp"].(bool); ok {
				dhcp = v
			}
		}
		if q, ok := body["quotas"].(map[string]any); ok {
			if v, ok := q["cores"].(float64); ok {
				cores = v
			}
		}
		mu.Unlock()
		jsonHandler(http.StatusOK, projectBody())(w, r)
	})
	srv.handle("DELETE /api/v3/projects/"+projectID, jsonHandler(http.StatusOK, map[string]any{}))

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resourceConfig := func(name string, version int, internetAccess, dhcp bool, cores int) string {
		return fmt.Sprintf(`
			resource "corellium_project" "test" {
				name = %q
				settings = {
					version         = %d
					internet_access = %t
					dhcp            = %t
				}
				quotas = {
					cores = %d
				}
				users = []
				teams = []
				keys  = []
			}`, name, version, internetAccess, dhcp, cores)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("test", 1, false, false, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_project.test", "id"),
					resource.TestCheckResourceAttr("corellium_project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + resourceConfig("test_update", 2, true, true, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_project.test", "name", "test_update"),
					resource.TestCheckResourceAttr("corellium_project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_project.test", "quotas.ram", "12288"),
				),
			},
		},
	})
}
