terraform {
  required_providers {
    corellium = {
      source  = "corellium/corellium"
      version = "~> 1.0.0"
    }
  }

  backend "s3" {}
}

provider "corellium" {
  # placeholder token - replace with real token or use env var CORELLIUM_TOKEN
  token = ""
}

resource "corellium_project" "example" {
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

resource "corellium_instance" "example" {
  name                   = "example"
  flavor                 = "iphone7plus"
  os_version             = "15.7.5"
  os_build               = "19H332"
  project                = corellium_project.example.id
  wait_for_ready         = true
  wait_for_ready_timeout = 500
}

resource "corellium_snapshot" "example" {
  name     = "example"
  instance = corellium_instance.example.id
}

