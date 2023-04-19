package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCorelliumV1ProjectResourceConfig() string {
	return fmt.Sprintf(
		`
		resource "corellium_v1project" "test" {
			name = "test"
			settings = {
				version = 1
				internet_access = false
				dhcp = false 
			}
			quotas = {
				cpus = 1
				cores = 1
				instances = 1
				ram = 1024 
			}
		}
		`,
	)
}

func TestAccCorelliumV1ProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// fa79475c-2703-4ccc-bc17-1876c2593ec6 is the "Terraform" project ID.
				Config: providerConfig + testAccCorelliumV1ProjectResourceConfig(),
				Check:  resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
				// resource.TestCheckResourceAttr("corellium_v1project.test", "type", "backup"),
				// resource.TestCheckResourceAttr("corellium_v1project.test", "project", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				),
			},
		},
	})
}
