package corellium

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1SofwareDataSource exercises the
// corellium_modelsoftware data source against a mock API server. This
// data source bypasses the generated client and builds its own request
// from the provider's configured host (see V1GetModelSoftwareManual in
// data_source_software_model.go), so this test also guards against that
// hand-rolled URL construction breaking again.
func TestUnitCorelliumV1SofwareDataSource(t *testing.T) {
	srv := newMockAPIServer(t)
	srv.handle("GET /api/v1/models/iphone7plus/software", jsonHandler(http.StatusOK, []map[string]any{
		{
			"api_version":    "",
			"android_flavor": "",
			"buildid":        "19H384",
			"filename":       "iPhone_7_Plus_15.7.5.ipsw",
			"id":             "abc",
			"md5sum":         "d41d8cd98f00b204e9800998ecf8427e",
			"orig_url":       "https://updates.cdn-apple.com/example.ipsw",
			"releasedate":    "2023-01-01",
			"sha1sum":        "da39a3ee5e6b4b0d3255bfef95601890afd80709",
			"sha256sum":      "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85",
			"size":           123456789,
			"uniqueid":       "unique-1",
			"uploaddate":     "2023-01-02",
			"url":            "https://example.com/iPhone_7_Plus_15.7.5.ipsw",
			"version":        "15.7.5",
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
                data "corellium_modelsoftware" "test" {
                    model = "iphone7plus"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.corellium_modelsoftware.test", "id"),
					resource.TestCheckResourceAttr("data.corellium_modelsoftware.test", "model_software.#", "1"),
					resource.TestCheckResourceAttr("data.corellium_modelsoftware.test", "model_software.0.build_id", "19H384"),
					resource.TestCheckResourceAttr("data.corellium_modelsoftware.test", "model_software.0.version", "15.7.5"),
					resource.TestCheckResourceAttr("data.corellium_modelsoftware.test", "model_software.0.size", "123456789"),
				),
			},
		},
	})
}
