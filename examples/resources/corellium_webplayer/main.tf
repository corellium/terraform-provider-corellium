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

resource "corellium_v1webplayer" "example" {
    project = ""
    instanceid = ""
    expiresinseconds = ""
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
