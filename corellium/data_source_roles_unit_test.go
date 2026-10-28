package corellium

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1RolesDataSource exercises the corellium_roles data
// source against a mock API server, with no TF_ACC, live token, or
// network egress beyond loopback.
func TestUnitCorelliumV1RolesDataSource(t *testing.T) {
	srv := newMockAPIServer(t)
	srv.handle("GET /api/v3/roles", jsonHandler(http.StatusOK, map[string]any{
		"results": []map[string]any{
			{"role": "admin", "project": "11111111-1111-1111-1111-111111111111", "user": "22222222-2222-2222-2222-222222222222"},
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
                data "corellium_roles" "test" {}
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_roles.test", "id"),
					resource.TestCheckResourceAttr("data.corellium_roles.test", "roles.#", "1"),
					resource.TestCheckResourceAttr("data.corellium_roles.test", "roles.0.role", "admin"),
					resource.TestCheckResourceAttr("data.corellium_roles.test", "roles.0.project", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("data.corellium_roles.test", "roles.0.user", "22222222-2222-2222-2222-222222222222"),
				),
			},
		},
	})
}
