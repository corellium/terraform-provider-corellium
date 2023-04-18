package corellium

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func preCheckAccCorelliumV1ImageResource() {
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

func testAccCorelliumV1ImageResourceConfig(name, kind, filename, encapsulated, project string) string {
	return fmt.Sprintf(
		`
		resource "corellium_v1image" "test" {
			name = "%s"
			type = "%s"
			filename = "%s"
			encapsulated = %s
			project = "%s"
		}
		`,
		name, kind, filename, encapsulated, project,
	)
}

func TestAccCorelliumV1ImageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		PreCheck:                 preCheckAccCorelliumV1ImageResource,
		Steps: []resource.TestStep{
			{
				// fa79475c-2703-4ccc-bc17-1876c2593ec6 is the "Terraform" project ID.
				Config: providerConfig + testAccCorelliumV1ImageResourceConfig("test", "backup", "/tmp/image.txt", "false", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1image.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "type", "backup"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "filename", "/tmp/image.txt"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "encapsulated", "false"),
					resource.TestCheckResourceAttr("corellium_v1image.test", "project", "fa79475c-2703-4ccc-bc17-1876c2593ec6"),
				),
			},
		},
	})
}
