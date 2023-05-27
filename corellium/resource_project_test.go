package corellium

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCorelliumV1ProjectResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test_update"
                    settings = {
                        version = 2
                        internet_access = true
                        dhcp = true
                    }
                    quotas = {
                        cores = 2
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test_update"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 2
                        internet_access = true
                        dhcp = true
                    }
                    quotas = {
                        cores = 2
                    }
                    users = [
                        {
                            id = "60d71152-8b86-4496-b27f-2e30f5bcc59f"
                            role = "admin"
                        }
                    ]
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.id", "60d71152-8b86-4496-b27f-2e30f5bcc59f"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 2
                        internet_access = true
                        dhcp = true
                    }
                    quotas = {
                        cores = 2
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 2
                        internet_access = true
                        dhcp = true
                    }
                    quotas = {
                        cores = 2
                    }
                    users = []
                    teams = [
                        {
                            id = "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"
                            role = "admin"
                        }
                    ]
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.id", "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 2
                        internet_access = true
                        dhcp = true
                    }
                    quotas = {
                        cores = 2
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 2
                        internet_access = true
                        dhcp = true
                    }
                    quotas = {
                        cores = 2
                    }
                    users = []
                    teams = []
					keys = [
					    {
							label = "test"
							kind = "ssh"
							key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEu0tbHR0DTV6wZqOcU/+xuNzyPsA7QzV9Eu5q2JbRVw"
						}
					]
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.0.kind", "ssh"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.0.key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEu0tbHR0DTV6wZqOcU/+xuNzyPsA7QzV9Eu5q2JbRVw"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_users(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = [
                        {
                            id = "60d71152-8b86-4496-b27f-2e30f5bcc59f"
                            role = "admin"
                        }
                    ]
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.id", "60d71152-8b86-4496-b27f-2e30f5bcc59f"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keyss.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keyss.#", "0"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_teams(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = [
                        {
                            id = "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"
                            role = "admin"
                        }
                    ]
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.id", "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_keys(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "corellium_v1project" "test" {
					name = "test"
					settings = {
						version = 1
						internet_access = false
						dhcp = false
					}
					quotas = {
						cores = 1
					}
					users = []
					teams = []
					keys = [
						{
							label = "test"
							kind = "ssh"
							key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFD33iT/L6sIb3kUWNMg2q9IbIF0DzksIRXJt4BbaP3K"
						}
					]
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.0.label", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.0.kind", "ssh"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.0.key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFD33iT/L6sIb3kUWNMg2q9IbIF0DzksIRXJt4BbaP3K"),
				),
			},
			{
				Config: providerConfig + `
				resource "corellium_v1project" "test" {
					name = "test"
					settings = {
						version = 1
						internet_access = false
						dhcp = false
					}
					quotas = {
						cores = 1
					}
					users = []
					teams = []
					keys = []
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_duplicated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = []
					keys = []
                }

                resource "corellium_v1project" "test_duplicated" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				ExpectError: regexp.MustCompile("A project with the name test already exists"),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_non_enterprise(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfigNonEnterprise + `
                resource "corellium_v1project" "test" {
                    name = "test"
                    settings = {
                        version = 1
                        internet_access = false
                        dhcp = false
                    }
                    quotas = {
                        cores = 1
                    }
                    users = []
                    teams = []
					keys = []
                }
                `,
				ExpectError: regexp.MustCompile("You don't have permission to create a project."),
			},
		},
	})
}
