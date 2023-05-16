# This is a demo Terraform configuration file for Corellium provider.
# This configuration file will add some users, create a team, a project, a instance and flavor.
terraform {
# The special terraform configuration block type is used to configure some behaviors of Terraform itself, such as
# requiring a minimum Terraform version to apply your configuration.
# Each terraform block can contain a number of settings related to Terraform's behavior. Within a terraform block, only
# constant values can be used; arguments may not refer to named objects such as resources, input variables, etc, and
# may not use any of the Terraform language built-in functions.
#
# https://developer.hashicorp.com/terraform/language/settings
  required_providers {
  # The required_providers block specifies all of the providers required by the current module, mapping each local
  # provider name to a source address and a version constraint.
  # Each Terraform module must declare which providers it requires, so that Terraform can install and use them.
  # Provider requirements are declared in a required_providers block.
  # A provider requirement consists of a local name, a source location, and a version constraint:
  #
  # https://developer.hashicorp.com/terraform/language/settings#specifying-provider-requirements
  # https://developer.hashicorp.com/terraform/language/providers/requirements#requiring-providers
  corellium = {
    source = "github.com/aimoda/corellium"
    # The global source address for the provider you intend to use, such as hashicorp/aws.
    #
    # https://developer.hashicorp.com/terraform/language/providers/requirements#source
    version = "~> 1.0.0"
    # A version constraint specifying which subset of available provider versions the module is compatible with.
    #
    # https://developer.hashicorp.com/terraform/language/providers/requirements#version

    # Each provider plugin has its own set of available versions, allowing the functionality of the provider to evolve
    # over time. Each provider dependency you declare should have a version constraint given in the version argument
    # so Terraform can select a single version per provider that all modules are compatible with.
    # The version argument is optional; if omitted, Terraform will accept any version of the provider as compatible.
    # However, we strongly recommend specifying a version constraint for every provider your module depends on.
    # To ensure Terraform always installs the same provider versions for a given configuration, you can use Terraform
    # CLI to create a dependency lock file and commit it to version control along with your configuration. If a lock
    # file is present, Terraform Cloud, CLI, and Enterprise will all obey it when installing providers.
    #
    # https://developer.hashicorp.com/terraform/language/providers/requirements#version-constraints
    }
  }
}

provider "corellium" {
# Providers allow Terraform to interact with cloud providers, SaaS providers, and other APIs.
# Some providers require you to configure them with endpoint URLs, cloud regions, or other settings before Terraform
# can use them.
# Additionally, all Terraform configurations must declare which providers they require so that Terraform can install
# and use them. The Provider Requirements page documents how to declare providers so Terraform can install them.
#
# https://developer.hashicorp.com/terraform/language/providers/configuration

# Provider configurations belong in the root module of a Terraform configuration. (Child modules receive their provider
# configurations from the root module; for more information, see the link below.)
#
# https://developer.hashicorp.com/terraform/language/providers/configuration#provider-configuration-1
  token = ""
  # Token is used to authenticate with the Corellium API. It can be provided via the CORELLIUM_TOKEN environment variable
  # or via the token argument. If both are provided, the token argument takes precedence.
}

locals {
# A local value assigns a name to an expression, so you can use the name multiple times within a module instead of 
# repeating the expression.
#
# https://developer.hashicorp.com/terraform/language/values/locals
  users = toset([
    "joao@testing.email.ai.moda",
    "maria@testing.email.ai.moda",
    "juan@testing.email.ai.moda",
    "luana@testing.email.ai.moda",
  ])
}

resource "random_string" "random" {
  length           = 32
  special          = true
  override_special = "/@£$"
}

resource "corellium_v1user" "student" {
# Resources are the most important element in the Terraform language. Each resource block describes one or more
# infrastructure objects, such as virtual networks, compute instances, or higher-level components such as DNS records.
#
# https://developer.hashicorp.com/terraform/language/resources/syntax

# Resource declarations can include a number of advanced features, but only a small subset are required for initial use.
# More advanced syntax features, such as single resource declarations that produce multiple similar remote objects,
# are described at the link below.
#
# https://developer.hashicorp.com/terraform/language/resources/syntax#resource-syntax
  for_each = local.users
  name = each.key
  label = each.key
  email = each.value
  password = random_string.random.result 
  administrator = false
}

resource "corellium_v1team" "demo" {
  label = "demo"
  users = [
    for user in corellium_v1user.student : {
      id = user.id
    }
  ]
}

# Create a project.
resource "corellium_v1project" "demo" {
  name = "demo"
  settings = {
    version = 1
    internet_access = false
    dhcp = false
  }
  quotas = {
    cores = 2
  }
  users = []
  teams = [
    {
      id = corellium_v1team.demo.id
      role = "admin"
    }
  ]
}
 
# Create a instance.
resource "corellium_v1instance" "student_instance" {
  name = "student_instance"
  # The flavor "samsung-galaxy-s-duos" is a Samsung Galaxy S Duos (GT-S7562) device and require 2 cores at least.
  flavor = "samsung-galaxy-s-duos"
  os = "13.0.0"
  project = corellium_v1project.demo.id
}
