package corellium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1SnapshotResource(t *testing.T) {
	t.Skip("Skipping snapshot resource tests")

	projectConfig := `
    resource "corellium_v1project" "test" {
        name = "test"
        settings = {
            version = 1
            internet_access = true
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

	instanceConfig := `
    resource "corellium_v1instance" "test" {
        name = "test"
        flavor = "iphone7plus"
        project = corellium_v1project.test.id
        os = "15.7.5"
        wait_for_ready = true
        wait_for_ready_timeout = 600
    }
    `

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + projectConfig + instanceConfig + `
                resource "corellium_v1snapshot" "test" {
                    name = "test"
                    instance = corellium_v1instance.test.id
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "name", "test"),
					resource.TestCheckResourceAttrSet("corellium_v1snapshot.test", "instance"),
				),
			},
			{
				Config: providerConfig + projectConfig + instanceConfig + `
                resource "corellium_v1snapshot" "test" {
                    name = "test_updat"
                    instance = corellium_v1instance.test.id
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "name", "test_updat"),
					resource.TestCheckResourceAttrSet("corellium_v1snapshot.test", "instance"),
				),
			},
		},
	})
}

func TestAccCorelliumV1SnapshotResource_non_enterprise(t *testing.T) {
	t.Skip("Skipping enterprise snapshot resource tests")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                data "corellium_v1projects" "test" {}

                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "samsung-galaxy-s-duos"
                    os = "13.0.0"
                    project = data.corellium_v1projects.test.projects[0].id
                }

                resource "corellium_v1snapshot" "test" {
                    name = "test"
                    instance = corellium_v1instance.test.id
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "name", "test"),
					resource.TestCheckResourceAttrSet("corellium_v1snapshot.test", "instance"),
				),
			},
		},
	})
}
