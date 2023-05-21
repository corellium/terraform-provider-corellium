package corellium

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1Users(t *testing.T) {
	resourceConfigCreateUser := func(administrator bool, label string, name string, email string, password string) string {
		return fmt.Sprintf(
			`
			resource "corellium_v1user" "test" {
				administrator = %t
				label = "%s"
				name = "%s"
				email = "%s"
				password = "%s"
			}
			`,
			administrator, label, name, email, password,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfigCreateUser(false, "ACC TEST ONE", "ACCTESTONENAME", "ACCTESTEMAILONE@ai.moda", "fdsahj29sd8dh@#$"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1user.test", "administrator", "false"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "label", "ACC TEST ONE"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "name", "ACCTESTONENAME"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "email", "ACCTESTEMAILONE@ai.moda"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "password", "fdsahj29sd8dh@#$"),
					resource.TestCheckNoResourceAttr("corellium_v1user.test", "ID"),
				),
			},
			{
				Config: providerConfig + resourceConfigCreateUser(false, "ACC TEST TWO", "ACCTESTTWONAME", "ACCTESTEMAILTWO@ai.moda", "fgdsg234e12eczx@#$"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1user.test", "administrator", "false"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "label", "ACC TEST TWO"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "name", "ACCTESTTWONAME"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "email", "ACCTESTEMAILTWO@ai.moda"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "password", "fgdsg234e12eczx@#$"),
					resource.TestCheckNoResourceAttr("corellium_v1user.test", "ID"),
				),
			},
			{
				Config: providerConfig + resourceConfigCreateUser(false, "ACC TEST THREE", "ACCTESTTHREENAME", "ACCTESTEMAILTHREE@ai.moda", "jbnmczxui8943211@#$"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1user.test", "administrator", "false"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "label", "ACC TEST THREE"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "name", "ACCTESTTHREENAME"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "email", "ACCTESTEMAILTHREE@ai.moda"),
					resource.TestCheckResourceAttr("corellium_v1user.test", "password", "jbnmczxui8943211@#$"),
					resource.TestCheckNoResourceAttr("corellium_v1user.test", "ID"),
				),
			},
		},
	})
}

func TestAccCorelliumV1Users_non_enterprise(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigNonEnterprise + fmt.Sprintf(`
                resource "corellium_v1user" "test" {
                    label = "test"
                    name = "test"
                    email = "test@testing.email.ai.moda"
                    password = "%s"
                    administrator = false
                }
                `, generatePassword(32, 4, 4, 4)),
				ExpectError: regexp.MustCompile("You do not have permission to create a user"),
			},
		},
	})
}
