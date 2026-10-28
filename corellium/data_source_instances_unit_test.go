package corellium

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1InstancesDataSource exercises the corellium_instances
// data source against a mock API server, with no TF_ACC, live token, or
// network egress beyond loopback.
func TestUnitCorelliumV1InstancesDataSource(t *testing.T) {
	srv := newMockAPIServer(t)
	srv.handle("GET /api/v3/instances", jsonHandler(http.StatusOK, map[string]any{
		"results": []map[string]any{
			{
				"id":           "55555555-5555-5555-5555-555555555555",
				"name":         "test",
				"key":          "abc123",
				"flavor":       "iphone7plus",
				"type":         "ios",
				"project":      "33333333-3333-3333-3333-333333333333",
				"state":        "on",
				"stateChanged": "2026-01-01T00:00:00Z",
				"startedAt":    "2026-01-01T00:00:00Z",
				"created":      "2026-01-01T00:00:00Z",
				"osVersion":    "15.7.5",
				"bootOptions":  map[string]any{},
				"netMon":       map[string]any{"hash": "", "info": "", "enabled": false},
			},
		},
		"count":   1,
		"total":   1,
		"limit":   50,
		"offset":  0,
		"hasMore": false,
	}))

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
                data "corellium_instances" "test" {}
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_instances.test", "id"),
					resource.TestCheckResourceAttr("data.corellium_instances.test", "instances.#", "1"),
					resource.TestCheckResourceAttr("data.corellium_instances.test", "instances.0.name", "test"),
					resource.TestCheckResourceAttr("data.corellium_instances.test", "instances.0.flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("data.corellium_instances.test", "instances.0.state", "on"),
				),
			},
		},
	})
}
