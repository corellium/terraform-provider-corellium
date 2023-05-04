package corellium

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1ImageResource(t *testing.T) {
	preCheck := func() {
		// It creates a file into the /tmp directory to be used as an image.
		// This is a workaround to avoid having to upload a real image to the
		// Corellium API to test the resource.
		if _, err := os.Stat("/tmp/image.txt"); os.IsNotExist(err) {
			f, err := os.Create("/tmp/image.txt")
			if err != nil {
				log.Println(err)
			}

			defer f.Close()
		}
	}

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
    }
    `

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 preCheck,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + projectConfig + `
                resource "corellium_v1image" "test" {
                    name = "test"
                    type = "backup"
                    filename = "/tmp/image.txt"
                    encapsulated = false
                    project = corellium_v1project.test.id
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1image.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "type", "backup"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "filename", "/tmp/image.txt"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "encapsulated", "false"),
					resource.TestCheckResourceAttrSet("corellium_v1image.test", "project"),
				),
			},
		},
	})
}
