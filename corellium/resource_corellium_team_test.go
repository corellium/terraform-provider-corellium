package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1TeamResource(t *testing.T) {
	// TODO: this could not be the best way to test this resource.

	resourceConfigWithoutUsers := func(label string) string {
		return fmt.Sprintf(
			`
		resource "corellium_v1team" "test" {
			label = "%s"
		}
		`, label,
		)
	}

	resourceConfigWithUsers := func(label string) string {
		// 60d71152-8b86-4496-b27f-2e30f5bcc59f is the ID of Henry Barreto.
		// TODO: create a user in the test and use its ID.
		return fmt.Sprintf(
			`
		resource "corellium_v1team" "test" {
			label = "%s"
			users = [
				{
					id = "60d71152-8b86-4496-b27f-2e30f5bcc59f"
				},
			]
		}
		`, label,
		)
	}

	resourceConfigEmptyUsers := func(label string) string {
		return fmt.Sprintf(
			`
		resource "corellium_v1team" "test" {
			label = "%s"
			users = []
		}
		`, label,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfigWithoutUsers("test_without_users"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_without_users"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "0"),
				),
			},
			{
				Config: providerConfig + resourceConfigWithoutUsers("test_without_users_label_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_without_users_label_updated"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "0"),
				),
			},
			{
				Config: providerConfig + resourceConfigWithUsers("test_with_users"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_with_users"),
					resource.TestCheckResourceAttrSet("corellium_v1team.test", "users.#"),
				),
			},
			{
				Config: providerConfig + resourceConfigWithUsers("test_with_users_label_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_with_users_label_updated"),
					resource.TestCheckResourceAttrSet("corellium_v1team.test", "users.#"),
				),
			},
			{
				Config: providerConfig + resourceConfigEmptyUsers("test_empty_users"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_empty_users"),
					resource.TestCheckResourceAttrSet("corellium_v1team.test", "users.#"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "0"),
				),
			},
			{
				Config: providerConfig + resourceConfigWithoutUsers("test_without_users_again"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_without_users_again"),
					resource.TestCheckNoResourceAttr("corellium_v1team.test", "users"),
				),
			},
		},
	})
}
