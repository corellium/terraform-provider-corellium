package corellium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1InstanceResource(t *testing.T) {
	projectConfig := `
    resource "corellium_v1project" "test" {
        name = "test"
        settings = {
            version = 1
            internet_access = false
            dhcp = false
        }
        quotas = {
            cores = 2
        }
        users = []
        teams = []
    }
    `

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + projectConfig + `
                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "iphone7plus"
                    project = corellium_v1project.test.id
                    os = "15.7.5"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "15.7.5"),
				),
			},
			{
				Config: providerConfig + projectConfig + `
                resource "corellium_v1instance" "test" {
                    name = "test_update"
                    flavor = "iphone7plus"
                    project = corellium_v1project.test.id
                    os = "15.7.5"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test_update"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "15.7.5"),
				),
			},
		},
	})
}
