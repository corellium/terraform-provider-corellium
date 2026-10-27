package corellium

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1SupportedModelsDataSource exercises the
// corellium_supportedmodels data source against a mock API server, with
// no TF_ACC, live token, or network egress beyond loopback.
func TestUnitCorelliumV1SupportedModelsDataSource(t *testing.T) {
	srv := newMockAPIServer(t)
	// The API client rejects a model missing any of its required
	// properties, so every one of them is present here even though the
	// data source only maps a handful.
	srv.handle("GET /api/v3/models", jsonHandler(http.StatusOK, []map[string]any{
		{
			"type":                        "iphone",
			"flavor":                      "iphone7plus",
			"description":                 "iPhone 7 Plus",
			"model":                       "N66AP",
			"hardwareModel":               "n66ap",
			"boardConfig":                 "n66ap",
			"platform":                    "t7000",
			"cpid":                        32777,
			"bdid":                        6,
			"vendor":                      "Apple",
			"hidden":                      false,
			"isBeta":                      false,
			"peripherals":                 false,
			"hasSoftware":                 true,
			"supportsRapidStart":          false,
			"defaultFirmware":             nil,
			"cores":                       2,
			"coresMin":                    2,
			"coresMax":                    2,
			"ram":                         2048,
			"ramMin":                      2048,
			"ramMax":                      2048,
			"dynamicHardwareChanges":      false,
			"partitions":                  []any{},
			"frame":                       nil,
			"productImage":                nil,
			"hasFrameVariant":             false,
			"showOptions":                 []string{},
			"stringUdid":                  false,
			"requiresPresetBootOptions":   false,
			"supportsPresetBootOptions":   false,
			"supportsCustomFirmware":      false,
			"supportsCustomFirmwareImage": false,
			"supportsCustomKernel":        false,
			"supportsCustomDevicetree":    false,
			"supportsCustomMmio":          false,
			"supportsFvp":                 false,
			"mast":                        false,
		},
	}))

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                data "corellium_supportedmodels" "test" {}
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_supportedmodels.test", "id"),
					resource.TestCheckResourceAttr("data.corellium_supportedmodels.test", "supported_models.#", "1"),
					resource.TestCheckResourceAttr("data.corellium_supportedmodels.test", "supported_models.0.type", "iphone"),
					resource.TestCheckResourceAttr("data.corellium_supportedmodels.test", "supported_models.0.model", "N66AP"),
					resource.TestCheckResourceAttr("data.corellium_supportedmodels.test", "supported_models.0.flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("data.corellium_supportedmodels.test", "supported_models.0.cp_id", "32777"),
					resource.TestCheckResourceAttr("data.corellium_supportedmodels.test", "supported_models.0.bd_id", "6"),
				),
			},
		},
	})
}
