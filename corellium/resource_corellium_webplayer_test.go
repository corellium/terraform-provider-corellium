package corellium

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1WebPlayer(t *testing.T) {
	resourceConfigCreateUser := func(project, instanceid string, expiresinseconds float32, apps, console, coretrace, devicecontrol, devicedelete, files, frida, images, messaging, netmon, network, portforwarding, profile, sensors, settings, snapshots, strace, system, connect bool) string {
		return fmt.Sprintf(
			`
			resource "corellium_v1webplayer" "test" {
				project = "%s"
				instanceid = "%s"
				expiresinseconds = %g
				features = {
				  apps = %t
				  console = %t
				  coretrace = %t
				  devicecontrol = %t
				  devicedelete = %t
				  files = %t
				  frida = %t
				  images = %t
				  messaging = %t
				  netmon = %t
				  network = %t
				  portforwarding = %t
				  profile = %t
				  sensors = %t
				  settings = %t
				  snapshots = %t
				  strace = %t
				  system = %t
				  connect = %t
				}
			  }
			`,
			project, instanceid, expiresinseconds, apps, console, coretrace, devicecontrol, devicedelete, files, frida, images, messaging, netmon, network, portforwarding, profile, sensors, settings, snapshots, strace, system, connect,
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfigCreateUser("3f054ee4-3732-40bf-ad95-f096f874db9f", "466b02f9-361d-40ef-a266-7dc6839b4522", 3600, true, true, true, true, false, true, false, false, true, true, true, false, true, true, true, true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "project", "3f054ee4-3732-40bf-ad95-f096f874db9f"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "instanceid", "466b02f9-361d-40ef-a266-7dc6839b4522"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "expiresinseconds", "3600"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.apps", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.console", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.coretrace", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.devicecontrol", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.devicedelete", "false"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.files", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.frida", "false"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.images", "false"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.messaging", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.netmon", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.network", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.portforwarding", "false"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.profile", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.sensors", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.settings", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.snapshots", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.strace", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.system", "true"),
					resource.TestCheckResourceAttr("corellium_v1webplayer.test", "features.connect", "true"),
					resource.TestCheckResourceAttrSet("corellium_v1webplayer.test", "token"),
					resource.TestCheckResourceAttrSet("corellium_v1webplayer.test", "identifier"),
					// resource.TestCheckResourceAttrSet("corellium_v1webplayer.test", "expiration"),
				),
			},
		},
	})
}
