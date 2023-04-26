package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Name          types.String `tfsdk:"name"`
// Label         types.String `tfsdk:"label"`
// Email         types.String `tfsdk:"email"`
// Password      types.String `tfsdk:"password"`
// Administrator types.Bool   `tfsdk:"administrator"`

func TestAccCorelliumV1Users(t *testing.T) {
	// TODO: this could not be the best way to test this resource.

	resourceConfigCreateUser := func(administrator bool, label string, name string, email string) string {
		return fmt.Sprintf(
			`
			resource "corellium_v1createuser" "test" {
				administrator = %t
				label = "%s"
				name = "%s"
				email = "%s"
			}
			`,
			administrator, label, name, email,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfigCreateUser(false, "testLabel", "test name", "test@test.com"),
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr("corellium_v1createuser.test", "administrator", "false"),
					// resource.TestCheckResourceAttr("corellium_v1createuser.test", "label", "testLabel"),
					// resource.TestCheckResourceAttr("corellium_v1createuser.test", "name", "test name"),
					// resource.TestCheckResourceAttr("corellium_v1createuser.test", "email", "test@test.com"),
					resource.TestCheckNoResourceAttr("corellium_v1createuser.test", "ID"),
					// resource.TestCheckResourceAttr("corellium_v1createuser.test", "password", "Ent12345@"),
				),
			},
		},
	})
}
