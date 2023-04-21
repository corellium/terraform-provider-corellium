package corellium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCorelliumV1SupportedModelsDataSourceConfig() string {
	return "\ndata \"corellium_v1supportedmodels\" \"test\" {}\n"
}

func TestAccCorelliumV1SupportedModelsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccCorelliumV1SupportedModelsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of supported models returned
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.#", "55"),
					// Verify the first model to ensure all attributes are set
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.type", "ios"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.name", "iphone14pm"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.flavor", "iphone14pm"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.description", "iPhone 14 Pro Max"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.model", "iPhone15,3"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.board_config", "d74ap"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.platform", "t8120"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.cp_id", "33056"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.bd_id", "14"),
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "supported_models.0.peripherals", "true"),
					// Verify placeholder id attribute since testing requires an id attribute to be set even though it is not used.
					resource.TestCheckResourceAttr("data.corellium_v1supportedmodels.test", "id", "placeholder"),
				),
			},
		},
	})
}
