<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

# Terraform Corellium Provider

![Acceptance Tests CI](https://github.com/aimoda/terraform-provider-corellium/actions/workflows/acceptance-tests.yml/badge.svg)
![Terraform Registry release CI](https://github.com/aimoda/terraform-provider-corellium/actions/workflows/terraform-release.yml/badge.svg)
![Examples QA CI](https://github.com/aimoda/terraform-provider-corellium/actions/workflows/examples-qa.yml/badge.svg)

The Corellium allows [Terraform](https://terraform.io) to manage [Corellium](https://www.corellium.com/?utm_source=github.com&utm_content=terraform-provider-corellium&utm_medium=github&utm_campaign=aimoda) resources.

- Examples can be found in the [examples](examples/) directory.
- Documentation can be found in the [docs](docs/) directory. At the same directory, you can find `demo.tf` file which contains all a small usage example.

_**Please note:** If you believe you have found a security issue in the Terraform Corellium Provider, please responsibly disclose it by contacting us at terraform-provider-corellium@security.email.ai.moda._

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x or higher

- [Go](https://golang.org/doc/install) 1.20.x (to build the provider plugin)

## Usage

This is a simple example of creating a project and multiple iOS versions at once.

```terraform
terraform {
  required_providers {
    corellium = {
      source  = "corellium/corellium"
      version = "0.0.2-alpha"
    }
  }
}

provider "corellium" {
  token = ""
}

resource "corellium_v1project" "backtesting" {
  name = "ios_backtesting"
  settings = {
    version         = 1
    internet_access = false
    dhcp            = false
  }
  quotas = {
    cores = 60
  }
  users = []
  teams = []
  keys = []
}

variable "ios_versions" {
  default = ["16.0", "16.0.2", "16.0.3", "16.1", "16.1.1", "16.1.2", "16.2", "16.3", "16.3.1", "16.4"]
  type    = set(string)
}

resource "corellium_v1instance" "test_instance" {
  for_each = var.ios_versions
  name     = "version_${each.key}"
  flavor   = "iphone8plus"
  project  = corellium_v1project.backtesting.id
  os       = each.key
}
```

Then, run:

```sh
CORELLIUM_API_TOKEN="YOUR.API_KEY_HERE" CORELLIUM_API_HOST="YOURDOMAIN.enterprise.corellium.com" terraform apply
```

To tear down everything that was created, run:
```sh
CORELLIUM_API_TOKEN="YOUR.API_KEY_HERE" CORELLIUM_API_HOST="YOURDOMAIN.enterprise.corellium.com" terraform destroy
```

<a href="https://www.ai.moda/en/?utm_source=github.com&utm_content=terraform-provider-corellium&utm_medium=github">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://terraform-provider-corellium.email.ai.moda/bimi/logo.svg?mode=dark">
    <source media="(prefers-color-scheme: light)" srcset="https://terraform-provider-corellium.email.ai.moda/bimi/logo.svg?mode=light">
    <img src="https://terraform-provider-corellium.email.ai.moda/bimi/logo.svg?mode=default" alt="ai.moda logo" title="ai.moda" align="right" height="50">
  </picture>
</a>

