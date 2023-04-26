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

resource "corellium_v1user" "example" {
	administrator = false
	label = "test label"
	name = "test name"
	email = "testemail@email.com"
  password = "testpassword"
}

