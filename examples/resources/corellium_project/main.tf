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

resource "corellium_v1project" "example" {
	name = "example"
	settings = {
		version = 1
		internet_access = false
		dhcp = false
	}
	quotas = {
		cores = 1
		instances = 2.5
		ram = 6144
	}
}

