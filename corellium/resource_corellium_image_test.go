package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCorelliumV1ImageResourceConfig(name, kind, project string) string {
	return fmt.Sprintf(
		`
		resource "corellium_v1image" "test" {
			name = "%s"
			type = "%s"
			project = "%s"
		}
		`,
		name, kind, project,
	)
}

func TestAccCorelliumV1ImageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// fa79475c-2703-4ccc-bc17-1876c2593ec6 is the "Terraform" project ID.
				Config: providerConfig + testAccCorelliumV1ImageResourceConfig("test", "backup", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1image.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "type", "backup"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "project", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				),
			},
		},
	})
}
