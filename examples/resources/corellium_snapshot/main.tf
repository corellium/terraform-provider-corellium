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

resource "corellium_v1snapshot" "example" {
	name = "example"
	instance = "1a968ea3-ac04-4201-9af0-25afdb80f4c4"
}

