package corellium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1InstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                data "corellium_v1instances" "test" {}
                `,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.corellium_v1projects.test", "instances.#"),
				),
			},
		},
	})
}

