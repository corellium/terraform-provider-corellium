package corellium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1InstanceResource_basic(t *testing.T) {
	projectConfig := `
    resource "corellium_v1project" "test" {
        name = "test"
        settings = {
            version = 1
            internet_access = false
            dhcp = false
        }
        quotas = {
            cores = 2
        }
        users = []
        teams = []
        keys  = []
    }
    `

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + projectConfig + `
                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "iphone7plus"
                    project = corellium_v1project.test.id
                    os = "15.7.5"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "15.7.5"),
				),
			},
			{
				Config: providerConfig + projectConfig + `
                resource "corellium_v1instance" "test" {
                    name = "test_update"
                    flavor = "iphone7plus"
                    project = corellium_v1project.test.id
                    os = "15.7.5"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test_update"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "iphone7plus"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "15.7.5"),
				),
			},
		},
	})
}

func TestAccCorelliumV1InstanceResource_wait_for_ready(t *testing.T) {
	projectConfig := `
    resource "corellium_v1project" "test" {
        name = "test"
        settings = {
            version = 1
            internet_access = false
            dhcp = false
        }
        quotas = {
            cores = 6
        }
        users = []
        teams = []
        keys  = []
    }
    `

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + projectConfig + `
                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "samsung-galaxy-s-duos"
                    project = corellium_v1project.test.id
                    os = "13.0.0"
                    wait_for_ready = true
                    wait_for_ready_timeout = 600
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "samsung-galaxy-s-duos"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "13.0.0"),
				),
			},
		},
	})
}

func TestAccCorelliumV1InstanceResource_default_project(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "samsung-galaxy-s-duos"
                    os = "13.0.0"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "samsung-galaxy-s-duos"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "13.0.0"),
				),
			},
		},
	})
}

func TestAccCorelliumV1InstanceResource_non_enterprise(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigNonEnterprise + `
                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "samsung-galaxy-s-duos"
                    os = "13.0.0"
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1instance.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "flavor", "samsung-galaxy-s-duos"),
					resource.TestCheckResourceAttr("corellium_v1instance.test", "os", "13.0.0"),
				),
			},
		},
	})
}
