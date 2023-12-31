terraform {
  required_providers {
    corellium = {
      source  = "github.com/aimoda/corellium"
      version = "~> 1.0.0"
    }
  }

  # backend "s3" {}
}

provider "corellium" {
  # placeholder token - replace with real token or use env var CORELLIUM_TOKEN
  token = ""
}

resource "corellium_v1project" "example" {
  name = "example"
  settings = {
    version         = 1
    internet_access = false
    dhcp            = false
  }
  quotas = {
    cores = 2
  }
  teams = []
  users = []
  keys  = []
}

resource "corellium_v1instance" "example" {
  project = corellium_v1project.example.id
  name    = "example"
  flavor  = "iphone7plus"
  os      = "15.7.5"

  provisioner "local-exec" {
    command = "uptime"
  }
}

