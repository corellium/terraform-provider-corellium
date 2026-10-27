package corellium

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1ProjectsDataSource exercises the corellium_projects
// data source against a mock API server, with no TF_ACC, live token, or
// network egress beyond loopback.
func TestUnitCorelliumV1ProjectsDataSource(t *testing.T) {
	const projectID = "33333333-3333-3333-3333-333333333333"

	srv := newMockAPIServer(t)
	srv.handle("GET /api/v3/projects", jsonHandler(http.StatusOK, map[string]any{
		"results": []map[string]any{
			{
				"id":   projectID,
				"name": "test",
				"settings": map[string]any{
					"version":        1,
					"internetAccess": false,
					"dhcp":           false,
				},
				"quotas": map[string]any{
					"cores":     2,
					"instances": 5,
					"ram":       6144,
				},
				"quotasUsed": map[string]any{},
				"keyCount":   1,
			},
		},
		"count":   1,
		"total":   1,
		"limit":   50,
		"offset":  0,
		"hasMore": false,
	}))
	srv.handle("GET /api/v3/projects/"+projectID+"/keys", jsonHandler(http.StatusOK, map[string]any{
		"results": []map[string]any{
			{
				"id":          "44444444-4444-4444-4444-444444444444",
				"project":     projectID,
				"kind":        "ssh",
				"key":         "ssh-rsa AAAA...",
				"fingerprint": "aa:bb:cc",
				"createdAt":   "2026-01-01T00:00:00Z",
				"updatedAt":   "2026-01-01T00:00:00Z",
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
                data "corellium_projects" "test" {}
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_projects.test", "id"),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.#", "1"),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.0.id", projectID),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.0.name", "test"),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.0.settings.version", "1"),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.0.quotas.cores", "2"),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.0.keys.#", "1"),
					resource.TestCheckResourceAttr("data.corellium_projects.test", "projects.0.keys.0.kind", "ssh"),
				),
			},
		},
	})
}
