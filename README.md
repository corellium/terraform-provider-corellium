<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

# Terraform Corellium Provider

![build](https://github.com/aimoda/terraform-provider-corellium/actions/workflows/devcontainer-build.yml/badge.svg)

The Corellium allows [Terraform](https://terraform.io) to manage [Corellium](https://www.corellium.com/) resources.

- Examples can be found in the [examples](examples/) directory.
- Documentation can be found in the [docs](docs/) directory. At the same directory, you can find `demo.tf` file which contains all a small usage example.

_**Please note:** If you believe you have found a security issue in the Terraform Corellium Provider, please responsibly disclose it by contacting us at terraform-provider-corellium@security.email.ai.moda._

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x or higher

- [Go](https://golang.org/doc/install) 1.20.x (to build the provider plugin)

## Usage

This is a simple example of creating a project and an instance for a user.

```terraform
terraform {
  required_providers {
    corellium = {
      source = "aimoda/corellium"
      version = "0.1.0"
    }
  }
}

provider "corellium" {
  token = ""
}

resource "corellium_v1project" "example" {
  name = "example"
  settings = {
    version = 1
    internet_access = true 
    dhcp = true
  }
  quotas = {
    cores = 2
  }
  users = [
    {
      id = "00000000-0000-0000-0000-000000000000"
    },
  ]
  teams = []
}

resource "corellium_v1instance" "example" {
  name = "student_instance"
  flavor = "samsung-galaxy-s-duos"
  project = corellium_v1project.example.id
  os = "13.0.0"
}
```