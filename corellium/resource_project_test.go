package corellium

import (
	"fmt"
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
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

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
                            id = corellium_v1user.test.id
                            role = "admin"
                        }
                    ]
                    teams = []
					keys = []
                }
                `, generatePassword(32, 4, 4, 4)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttrSet("corellium_v1project.test", "users.0.id"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

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
                `, generatePassword(32, 4, 4, 4)),
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
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

				resource "corellium_v1team" "test" {
					label = "test"
					users = [
						{
							id = corellium_v1user.test.id
						},
					]
				}

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
                            id = corellium_v1team.test.id
                            role = "admin"
                        }
                    ]
					keys = []
                }
                `, generatePassword(32, 4, 4, 4)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttrSet("corellium_v1project.test", "teams.0.id"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

				resource "corellium_v1team" "test" {
					label = "test"
					users = [
						{
							id = corellium_v1user.test.id
						},
					]
				}

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
                `, generatePassword(32, 4, 4, 4)),
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
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

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
                            id = corellium_v1user.test.id
                            role = "admin"
                        }
                    ]
                    teams = []
					keys = []
                }
                `, generatePassword(32, 4, 4, 4)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttrSet("corellium_v1project.test", "users.0.id"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keyss.#", "0"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

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
                `, generatePassword(32, 4, 4, 4)),
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
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

				resource "corellium_v1team" "test" {
					label = "test"
					users = [
						{
							id = corellium_v1user.test.id
						},
					]
				}

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
                            id = corellium_v1team.test.id
                            role = "admin"
                        }
                    ]
					keys = []
                }
                `, generatePassword(32, 4, 4, 4)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
					resource.TestCheckResourceAttrSet("corellium_v1project.test", "teams.0.id"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.role", "admin"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "keys.#", "0"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
				resource "corellium_v1user" "test" {
					label = "test"
					name = "test"
					email = "testing@testing.ai.moda"
					password = "%s"
					administrator = true
				}

				resource "corellium_v1team" "test" {
					label = "test"
					users = [
						{
							id = corellium_v1user.test.id
						},
					]
				}

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
                `, generatePassword(32, 4, 4, 4)),
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
