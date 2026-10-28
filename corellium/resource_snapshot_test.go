package corellium

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1SnapshotResource(t *testing.T) {
	t.Skip("Skipping snapshot resource tests")

	projectConfig := `
    resource "corellium_project" "test" {
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
    resource "corellium_instance" "test" {
        name = "test"
        flavor = "iphone7plus"
        project = corellium_project.test.id
        os_version = "15.7.5"
        os_build = "19H332"
        wait_for_ready = true
        wait_for_ready_timeout = 600
    }
    `

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + projectConfig + instanceConfig + `
                resource "corellium_snapshot" "test" {
                    name = "test"
                    instance = corellium_instance.test.id
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_snapshot.test", "name", "test"),
					resource.TestCheckResourceAttrSet("corellium_snapshot.test", "instance"),
				),
			},
			{
				Config: providerConfig + projectConfig + instanceConfig + `
                resource "corellium_snapshot" "test" {
                    name = "test_updat"
                    instance = corellium_instance.test.id
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_snapshot.test", "name", "test_updat"),
					resource.TestCheckResourceAttrSet("corellium_snapshot.test", "instance"),
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
                data "corellium_projects" "test" {}

                resource "corellium_instance" "test" {
                    name = "test"
                    flavor = "samsung-galaxy-s-duos"
                    os_version = "13.0.0"
                    os_build = ""
                    project = data.corellium_projects.test.projects[0].id
                }

                resource "corellium_snapshot" "test" {
                    name = "test"
                    instance = corellium_instance.test.id
                }`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_snapshot.test", "name", "test"),
					resource.TestCheckResourceAttrSet("corellium_snapshot.test", "instance"),
				),
			},
		},
	})
}
