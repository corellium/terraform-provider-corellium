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
					// Verify the first model to ensure all attributes are set
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model", "iPhone15,3"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.version", "16.0"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.build_id", "20A362"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.sha256_sum", "4f12cc262aa87647bebc83a6ca0ae29bd11dff2ad73812cc45d920be4caaa584"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.sha1_sum", ""),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.md5_sum", "fa07e00024db4a50a132cd4ba7575c10-788"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.size", "6609575844"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.unique_id", ""),
					// resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.metadata", ""),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.release_date", "2022-09-12T17:01:06Z"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.upload_date", ""),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.url", "https://updates.cdn-apple.com/2022FallFCS/fullrestores/012-65663/2051812C-0862-4EA6-A896-365466C2DBAD/iPhone15,3_16.0_20A362_Restore.ipsw"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.orig_url", "https://updates.cdn-apple.com/2022FallFCS/fullrestores/012-65663/2051812C-0862-4EA6-A896-365466C2DBAD/iPhone15,3_16.0_20A362_Restore.ipsw"),
					resource.TestCheckResourceAttr("data.corellium_v1modelsoftware.test", "model_software.0.filename", "iPhone15,3_16.0_20A362_Restore.ipsw"),
					resource.TestCheckResourceAttrSet("data.corellium_v1modelsoftware.test", "id"),
				),
			},
		},
	})
}
