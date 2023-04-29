package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1InstanceResource(t *testing.T) {
	resourceConfig := func(name, flavor, project, os string) string {
		return fmt.Sprintf(
			`
		resource "corellium_v1instance" "test" {
            name = "%s"
			flavor = "%s"
			project = "%s"
			os = "%s"
		}
		`,
			name, flavor, project, os,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("test", "iphone7plus", "fa79475c-2703-4ccc-bc17-1876c2593ec6", "15.7.5"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "project", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "15.7.5"),
				),
			},
			{
				Config: providerConfig + resourceConfig("test_updated", "iphone7plus", "fa79475c-2703-4ccc-bc17-1876c2593ec6", "15.7.5"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test_updated"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "project", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "15.7.5"),
				),
			},
		},
	})
}
