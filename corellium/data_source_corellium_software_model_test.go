package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCorelliumV1SoftwareDataSourceConfig(model string) string {
	return fmt.Sprintf(
		`
		data "corellium_v1modelsoftware" "test" {
			model = "%s"
		}
		`,
		model,
	)
}

func TestAccCorelliumV1SofwareDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCorelliumV1SoftwareDataSourceConfig("iPhone15,3"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1modelsoftware.test", "model", "iPhone15,3"),
				),
			},
		},
	})
}
