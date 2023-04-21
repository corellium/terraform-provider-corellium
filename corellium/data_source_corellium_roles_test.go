package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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
		PreCheck:                 preCheckAccCorelliumV1ImageResource,
		Steps: []resource.TestStep{
			{
				// Check that the data source can be read without any arguments.
				Config: providerConfig + testAccCorelliumV1RolesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.#"),
				),
			},
			{
				// Check that the data source can be read with a project argument.
				// fa79475c-2703-4ccc-bc17-1876c2593ec6 is the "Terraform" project ID.
				Config: providerConfig + testAccCorelliumV1RolesDataSourceConfigWithProject("fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.#"),
					resource.TestCheckResourceAttr("data.corellium_v1roles.test", "roles.0.project", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				),
			},
			{
				// Check that the data source can be read with an invalid project argument.
				// In this case, the data source should return a list with all roles from all projects and users.
				Config: providerConfig + testAccCorelliumV1RolesDataSourceConfigWithProject("invalid or unknown"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_v1roles.test", "roles.#"),
				),
			},
		},
	})
}
