terraform {
  required_providers {
    corellium = {
      source = "github.com/aimoda/corellium"
      version = "~> 1.0.0"
    }
  }
}

provider "corellium" {
  # placeholder token - replace with real token or use env var CORELLIUM_TOKEN
  token = ""
}

resource "random_string" "random" {
  length           = 32
  special          = true
  override_special = "/@£$"
}

resource "corellium_v1user" "example" {
  label = "example"
  name = "example"
  email = "example@testing.email.ai.moda"
  password = random_string.random.result
  administrator = false
}

