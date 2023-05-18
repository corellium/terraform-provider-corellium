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
                    name = "test_create"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test_create"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "6144"),
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
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test_update_add_users"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test_update_add_users"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.id", "60d71152-8b86-4496-b27f-2e30f5bcc59f"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.0.role", "admin"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test_update_remove_users"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test_update_remove_users"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "users.#", "0"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test_update_add_teams"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test_update_add_teams"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.id", "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.0.role", "admin"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test" {
                    name = "test_update_remove_teams"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test", "name", "test_update_remove_teams"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.version", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.internet_access", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "settings.dhcp", "true"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.cores", "2"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.instances", "5"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "quotas.ram", "12288"),
					resource.TestCheckResourceAttr("corellium_v1project.test", "teams.#", "0"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_add_users_on_creation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test_create_with_users" {
                    name = "test_create_with_users"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "name", "test_create_with_users"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "users.0.id", "60d71152-8b86-4496-b27f-2e30f5bcc59f"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_users", "users.0.role", "admin"),
				),
			},
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test_create_with_teams" {
                    name = "test_create_with_teams"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "name", "test_create_with_teams"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "teams.0.id", "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "teams.0.role", "admin"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_add_teams_on_creation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
                resource "corellium_v1project" "test_create_with_teams" {
                    name = "test_create_with_teams"
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
                }
                `,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "name", "test_create_with_teams"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "settings.version", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "settings.internet_access", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "settings.dhcp", "false"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "quotas.cores", "1"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "quotas.instances", "2.5"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "quotas.ram", "6144"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "teams.0.id", "d1c3d32d-b46e-4ba3-8f31-2069fc5e80bc"),
					resource.TestCheckResourceAttr("corellium_v1project.test_create_with_teams", "teams.0.role", "admin"),
				),
			},
		},
	})
}

func TestAccCorelliumV1ProjectResource_project_name_duplicated(t *testing.T) {
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
                }
                `,
				ExpectError: regexp.MustCompile("You don't have permission to create a project."),
			},
		},
	})
}
