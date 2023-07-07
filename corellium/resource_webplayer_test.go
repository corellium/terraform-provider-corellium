package corellium

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1WebPlayer(t *testing.T) {
	t.Skip("The current enterprise account doesn't have WebPlayer enabled. It'll be fixed in the future.")

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
        wait_for_ready = false 
    }
    `
	webplayerConfig := func(project, instanceid string, expiresinseconds float32, apps, console, coretrace, devicecontrol, devicedelete, files, frida, images, messaging, netmon, network, portforwarding, profile, sensors, settings, snapshots, strace, system, connect bool) string {
		return fmt.Sprintf(
			`
			resource "corellium_v1webplayer" "test" {
				project = %s
				instanceid = %s
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
				Config: providerConfig + projectConfig + instanceConfig + webplayerConfig("corellium_v1project.test.id", "corellium_v1instance.test.id", 3600, true, true, true, true, false, true, false, false, true, true, true, false, true, true, true, true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr("corellium_v1webplayer.test", "project", "3f054ee4-3732-40bf-ad95-f096f874db9f"),
					// resource.TestCheckResourceAttr("corellium_v1webplayer.test", "instanceid", "466b02f9-361d-40ef-a266-7dc6839b4522"),
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

func TestAccCorelliumV1WebPlayer_non_enterprise(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigNonEnterprise + `
                data "corellium_v1projects" "test" {}

                resource "corellium_v1instance" "test" {
                    name = "test"
                    flavor = "samsung-galaxy-s-duos"
                    os = "13.0.0"
                    project = data.corellium_v1projects.test.projects[0].id
                }

                resource "corellium_v1webplayer" "test" {
                    project = data.corellium_v1projects.test.projects[0].id
                    instanceid = corellium_v1instance.test.id
                    expiresinseconds = 0
                    features = {
                      apps = false
                      console = false
                      coretrace = false
                      devicecontrol = false
                      devicedelete = false
                      files = false
                      frida = false
                      images = false
                      messaging = false
                      netmon = false
                      network = false
                      portforwarding = false
                      profile = false
                      sensors = false
                      settings = false
                      snapshots = false
                      strace = false
                      system = false
                      connect = false
                    }
                }
                `,
				ExpectError: regexp.MustCompile("Error creating a web player session"),
			},
		},
	})
}
