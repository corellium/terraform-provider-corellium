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

resource "corellium_v1instance" "example" {
    name = "example"
    flavor = "iphone7plus"
    project = "fa79475c-2703-4ccc-bc17-1876c2593ec6"
    os = "15.7.5"
}

