package corellium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1InstanceResource exercises the corellium_instance
// resource's Create/Read/Update lifecycle against a mock API server, with
// wait_for_ready disabled to avoid exercising the retry.StateChangeConf
// polling loop (that loop is a thin wrapper around repeated GETs and isn't
// this resource's core logic). No TF_ACC, live token, or network egress
// beyond loopback.
func TestUnitCorelliumV1InstanceResource(t *testing.T) {
	const instanceID = "55555555-5555-5555-5555-555555555555"
	const projectID = "33333333-3333-3333-3333-333333333333"

	var mu sync.Mutex
	name := "test"
	deleted := false

	instanceBody := func() map[string]any {
		mu.Lock()
		defer mu.Unlock()
		return map[string]any{
			"id":           instanceID,
			"name":         name,
			"key":          "abc123",
			"flavor":       "iphone7plus",
			"type":         "ios",
			"project":      projectID,
			"state":        "on",
			"stateChanged": "2026-01-01T00:00:00Z",
			"startedAt":    "2026-01-01T00:00:00Z",
			"created":      "2026-01-01T00:00:00Z",
			"osVersion":    "15.7.5",
			"bootOptions":  map[string]any{},
		}
	}

	srv := newMockAPIServer(t)
	srv.handle("POST /api/v3/instances", jsonHandler(http.StatusOK, map[string]any{"id": instanceID, "state": "creating"}))
	srv.handle("GET /api/v3/instances/"+instanceID, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		gone := deleted
		mu.Unlock()
		if gone {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		jsonHandler(http.StatusOK, instanceBody())(w, r)
	})
	srv.handle("PATCH /api/v3/instances/"+instanceID, func(w http.ResponseWriter, r *http.Request) {
		var patch map[string]any
		_ = json.NewDecoder(r.Body).Decode(&patch)
		if n, ok := patch["name"].(string); ok {
			mu.Lock()
			name = n
			mu.Unlock()
		}
		jsonHandler(http.StatusOK, instanceBody())(w, r)
	})
	srv.handle("DELETE /api/v3/instances/"+instanceID, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		deleted = true
		mu.Unlock()
		jsonHandler(http.StatusOK, map[string]any{})(w, r)
	})

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resourceConfig := func(name string) string {
		return fmt.Sprintf(`
			resource "corellium_instance" "test" {
				name           = %q
				flavor         = "iphone7plus"
				project        = %q
				os_version     = "15.7.5"
				os_build       = ""
				wait_for_ready = false
			}`, name, projectID)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_instance.test", "id"),
					resource.TestCheckResourceAttr("corellium_instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_instance.test", "os_version", "15.7.5"),
					resource.TestCheckResourceAttr("corellium_instance.test", "state", "on"),
				),
			},
			{
				Config: providerConfig + resourceConfig("test_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_instance.test", "name", "test_updated"),
				),
			},
		},
	})
}
