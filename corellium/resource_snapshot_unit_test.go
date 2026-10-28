package corellium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1SnapshotResource exercises the corellium_snapshot
// resource's full CRUD lifecycle against a mock API server. No TF_ACC,
// live token, or network egress beyond loopback.
func TestUnitCorelliumV1SnapshotResource(t *testing.T) {
	const snapshotID = "99999999-9999-9999-9999-999999999999"
	const instanceID = "55555555-5555-5555-5555-555555555555"

	var mu sync.Mutex
	name := "test"

	snapshotBody := func() map[string]any {
		mu.Lock()
		defer mu.Unlock()
		return map[string]any{
			"id":       snapshotID,
			"name":     name,
			"instance": instanceID,
			"status":   map[string]any{"task": "none", "created": true},
			"date":     1735689600000.0,
			"fresh":    false,
			"live":     false,
			"local":    true,
			"sharing":  map[string]any{"sharingType": "publicLink"},
		}
	}

	srv := newMockAPIServer(t)
	srv.handle("POST /api/v3/instances/"+instanceID+"/snapshots", func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(http.StatusOK, snapshotBody())(w, r)
	})
	srv.handle("GET /api/v3/snapshots/"+snapshotID, func(w http.ResponseWriter, r *http.Request) {
		jsonHandler(http.StatusOK, snapshotBody())(w, r)
	})
	srv.handle("PATCH /api/v3/instances/"+instanceID+"/snapshots/"+snapshotID, func(w http.ResponseWriter, r *http.Request) {
		var patch map[string]any
		_ = json.NewDecoder(r.Body).Decode(&patch)
		if n, ok := patch["name"].(string); ok {
			mu.Lock()
			name = n
			mu.Unlock()
		}
		jsonHandler(http.StatusOK, snapshotBody())(w, r)
	})
	srv.handle("DELETE /api/v3/snapshots/"+snapshotID, jsonHandler(http.StatusOK, map[string]any{}))

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resourceConfig := func(name string) string {
		return fmt.Sprintf(`
			resource "corellium_snapshot" "test" {
				name     = %q
				instance = %q
			}`, name, instanceID)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_snapshot.test", "id"),
					resource.TestCheckResourceAttr("corellium_snapshot.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_snapshot.test", "instance", instanceID),
					resource.TestCheckResourceAttr("corellium_snapshot.test", "status.task", "none"),
				),
			},
			{
				Config: providerConfig + resourceConfig("test_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_snapshot.test", "name", "test_updated"),
				),
			},
		},
	})
}
