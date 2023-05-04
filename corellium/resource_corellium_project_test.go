package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1ProjectResource(t *testing.T) {
	resourceConfig := func(name, version, internet_access, dhcp, cores string) string {
		return fmt.Sprintf(
			`
		resource "corellium_v1project" "test" {
			name = "%s"
			settings = {
				version = %s
				internet_access = %s 
				dhcp = %s 
			}
			quotas = {
				cores = %s
			}
		}
		`, name, version, internet_access, dhcp, cores,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("test", "1", "false", "false", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
				),
			},
			{
				Config: providerConfig + resourceConfig("test2", "2", "true", "true", "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
				),
			},
		},
	})
}
