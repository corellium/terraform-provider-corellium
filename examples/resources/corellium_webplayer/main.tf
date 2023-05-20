terraform {
  required_providers {
    corellium = {
      source = "github.com/aimoda/corellium"
      version = "~> 1.0.0"
    }
  }

  backend "s3" {}
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
      cores = 2
  }
  users = []
  teams = []
}

resource "corellium_v1instance" "example" {
  name = "example"
  flavor = "ranchu"
  project = corellium_v1project.example.id
  os = "7.1.2"
}

resource "corellium_v1webplayer" "example" {
  project = corellium_v1project.example.id
  instanceid = corellium_v1instance.example.id
  expiresinseconds = "1800"
  features = {
    apps = false
    console = false
    coretrace = false
    devicecontrol = false
    devicedelete = false
    files = false
    frida = false
    images = false
    messaging = false
    netmon = false
    network = false
    portforwarding = false
    profile = false
    sensors = false
    settings = false
    snapshots = false
    strace = false
    system = false
    connect = false
}
}
