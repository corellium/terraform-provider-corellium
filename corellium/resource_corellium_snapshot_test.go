package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1SnapshotResource(t *testing.T) {
	resourceConfig := func(name, instance string) string {
		// 1a968ea3-ac04-4201-9af0-25afdb80f4c4 is the ID of the instance "phase-partridge".
		return fmt.Sprintf(
			`
		resource "corellium_v1snapshot" "test" {
			name = "%s"
			instance = "%s"
		}
		`, name, instance,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig("test", "1a968ea3-ac04-4201-9af0-25afdb80f4c4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "instance", "1a968ea3-ac04-4201-9af0-25afdb80f4c4"),
				),
			},
			{
				Config: providerConfig + resourceConfig("test_updated", "1a968ea3-ac04-4201-9af0-25afdb80f4c4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "name", "test_updated"),
					resource.TestCheckResourceAttr("corellium_v1snapshot.test", "instance", "1a968ea3-ac04-4201-9af0-25afdb80f4c4"),
				),
			},
		},
	})
}
