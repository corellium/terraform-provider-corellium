package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1ReadyDataSource is the harness smoke test: it proves
// the provider can be pointed at a local mock HTTPS server (via the "host"
// configuration attribute) end-to-end, with no TF_ACC, no live API token,
// and no network egress beyond loopback.
func TestUnitCorelliumV1ReadyDataSource(t *testing.T) {
	srv := newMockAPIServer(t)

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
                data "corellium_ready" "test" {}
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_ready.test", "id"),
					resource.TestCheckResourceAttrSet("data.corellium_ready.test", "status"),
				),
			},
		},
	})
}
