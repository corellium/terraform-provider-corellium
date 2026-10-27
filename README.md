<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

# Terraform Corellium Provider

![Acceptance Tests CI](https://github.com/corellium/terraform-provider-corellium/actions/workflows/acceptance-tests.yml/badge.svg)
![Terraform Registry release CI](https://github.com/corellium/terraform-provider-corellium/actions/workflows/terraform-release.yml/badge.svg)
![Examples QA CI](https://github.com/corellium/terraform-provider-corellium/actions/workflows/examples-qa.yml/badge.svg)

The Corellium allows [Terraform](https://terraform.io) to manage [Corellium](https://www.corellium.com/?utm_source=github.com&utm_content=terraform-provider-corellium&utm_medium=github) resources.

- Examples can be found in the [examples](examples/) directory.
- Documentation can be found in the [docs](docs/) directory. At the same directory, you can find `demo.tf` file which contains all a small usage example.

_**Please note:** If you believe you have found a security issue in the Terraform Corellium Provider, please responsibly disclose it by contacting us at terraform-provider-corellium@security.email.ai.moda._

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.13.x or higher

- [Go](https://golang.org/doc/install) 1.25.x (if you want to build the provider plugin)

## Quick Start

This is a simple example of creating a project with multiple devices at once.

### Set the provider version in the 'required_providers' block

```
terraform {
  required_providers {
    corellium = {
      source  = "corellium/corellium"
      version = "<provider version here>"
    }
  }
}
```

### Set the provider host and API token

```
provider "corellium" {
  token = "<API key here>"
  host = "<host machine here>"
}
```

### Alternatively, you can set host and token as environment variables. This would help avoid pushing the token to source control by mistake

```
provider "corellium" {

}
```

And set the environment variable in your termianl

Linux:

```
export CORELLIUM_API_TOKEN="<your token here>"
export CORELLIUM_API_HOST="<host machine here>"
```

Windows:

```
set CORELLIUM_API_TOKEN="<your token here>"
set CORELLIUM_API_HOST="<host machine here>"
```



### Define a project. Set the number of cores it uses and other settings

```
resource "corellium_project" "proj" {
  name = "demo_project"
  settings = {
    version         = 1
    internet_access = false
    dhcp            = true
  }
  quotas = {
    cores = 20
  }
  users = []
  teams = []
  keys = []
}
```

### Define the devices that are part of the project

```
resource "corellium_instance" "android" {
   project = corellium_project.proj.id
   name    = "Demo Android"
   flavor  = "ranchu"
   os_version      = "16.0.0"
   os_build = "r2 userdebug"
}
 
resource "corellium_instance" "ios" {
   name     = "Demo iOS2"
   flavor   = "iphone17e"
   project  = corellium_project.proj.id
   os_version       = "26.5.2"
   os_build = "23F84"
}
```

Then, run:

```
terraform init
```

This will initialize the provider. To see the plan summary, run:

```
terraform plan
```

If everything looks good, apply the configuration: (type 'yes' to approve when prompt)

```
terraform apply 
```

To tear down everything that was created, run:

```
terraform destroy 
```

<a href="https://www.ai.moda/en/?utm_source=github.com&utm_content=terraform-provider-corellium&utm_medium=github">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://terraform-provider-corellium.email.ai.moda/bimi/logo.svg?mode=dark">
    <source media="(prefers-color-scheme: light)" srcset="https://terraform-provider-corellium.email.ai.moda/bimi/logo.svg?mode=light">
    <img src="https://terraform-provider-corellium.email.ai.moda/bimi/logo.svg?mode=default" alt="ai.moda logo" title="ai.moda" align="right" height="50">
  </picture>
</a>

