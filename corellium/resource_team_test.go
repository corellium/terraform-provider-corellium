package corellium

import (
	"fmt"
	"regexp"
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

	resourceConfigOneUser := func(label string) string {
		return fmt.Sprintf(
			`
		resource "corellium_v1user" "test" {
			label = "test"
			name = "test"
			email = "testing@testing.ai.moda"
			password = "%s"
			administrator = true
		}

		resource "corellium_v1team" "test" {
			label = "%s"
			users = [
				{
					id = corellium_v1user.test.id
				},
			]
		}
		`, generatePassword(32, 4, 4, 4), label,
		)
	}

	resourceConfigTwoUsers := func(label string) string {
		// 60d71152-8b86-4496-b27f-2e30f5bcc59f is the ID of Henry Barreto.
		// TODO: create a user in the test and use its ID.
		return fmt.Sprintf(
			`
		resource "corellium_v1user" "test" {
			label = "test"
			name = "test"
			email = "testing@testing.ai.moda"
			password = "%s"
			administrator = true
		}

		resource "corellium_v1user" "user" {
			label = "user"
			name = "user"
			email = "user@testing.ai.moda"
			password = "%s"
			administrator = false
		}

		resource "corellium_v1team" "test" {
			label = "%s"
			users = [
				{
					id = corellium_v1user.test.id
				},
                {
                    id = corellium_v1user.user.id
                }
			]
		}
		`, generatePassword(32, 4, 4, 4), generatePassword(32, 4, 4, 4), label,
		)
	}

	resourceConfigEmptyUsers := func(label string) string {
		return fmt.Sprintf(
			`
		resource "corellium_v1user" "test" {
			label = "test"
			name = "test"
			email = "testing@testing.ai.moda"
			password = "%s"
			administrator = true
		}

		resource "corellium_v1user" "user" {
			label = "user"
			name = "user"
			email = "user@testing.ai.moda"
			password = "%s"
			administrator = false
		}

		resource "corellium_v1team" "test" {
			label = "%s"
			users = []
		}
		`, generatePassword(32, 4, 4, 4), generatePassword(32, 4, 4, 4), label,
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
				Config: providerConfig + resourceConfigOneUser("test_with_one_user"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_with_one_user"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "1"),
				),
			},
			{
				Config: providerConfig + resourceConfigTwoUsers("test_with_two_users"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_with_two_users"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "2"),
				),
			},
			{
				Config: providerConfig + resourceConfigTwoUsers("test_with_users_label_updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_with_users_label_updated"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "2"),
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
				Config: providerConfig + resourceConfigTwoUsers("test_with_users_label_updated_again"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1team.test", "label", "test_with_users_label_updated_again"),
					resource.TestCheckResourceAttr("corellium_v1team.test", "users.#", "2"),
				),
			},
		},
	})
}

func TestAccCorelliumV1TeamResource_non_enterprise(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigNonEnterprise + `
                resource "corellium_v1team" "test" {
                    label = "test"
                    users = []
                }
                `,
				ExpectError: regexp.MustCompile("You don't have permissions to create a team"),
			},
		},
	})
}
