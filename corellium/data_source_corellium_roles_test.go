package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TODO: remove these help functions inserting into the acc test.
func testAccCorelliumV1RolesDataSourceConfig() string {
	return fmt.Sprintf(
		`
		data "corellium_v1roles" "test" { }

		output "roles" {
			value = data.corellium_v1roles.test.roles
		}
		`,
	)
}

func testAccCorelliumV1RolesDataSourceConfigWithProject(project string) string {
	return fmt.Sprintf(
		`
		data "corellium_v1roles" "test" {
			project = "%s"
		}

		output "roles" {
			value = data.corellium_v1roles.test.roles
		}
		`, project,
	)
}

func TestAccCorelliumV1RolesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCorelliumV1RolesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.#"),
				),
			},
			{
				Config: providerConfig + testAccCorelliumV1RolesDataSourceConfigWithProject("fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.#"),
                    resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.0.project"),
				),
			},
			{
				Config: providerConfig + testAccCorelliumV1RolesDataSourceConfigWithProject("invalid or unknown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.#"),
				),
			},
		},
	})
}
