terraform {
  required_providers {
    corellium = {
      source  = "github.com/aimoda/corellium"
      version = "~> 1.0.0"
    }
  }

  backend "s3" {}
}

provider "corellium" {
  # placeholder token - replace with real token or use env var CORELLIUM_TOKEN
  token = ""
}

resource "tls_private_key" "example" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "corellium_v1project" "example" {
  name = "example"
  settings = {
    version         = 1
    internet_access = true
    dhcp            = true
  }
  quotas = {
    cores = 4
  }
  users = []
  teams = []
  keys = [
    {
      label = "example"
      kind  = "ssh"
      key   = tls_private_key.example.public_key_openssh
    }
  ]
}

resource "corellium_v1instance" "example" {
  project                = corellium_v1project.example.id
  name                   = "example"
  flavor                 = "iphone6splus"
  os                     = "13.6.1"
  wait_for_ready         = true
  wait_for_ready_timeout = 900

  connection {
    # corellium instances supports ssh and adb connections, but terraform only supports ssh connections, so we can only
    # use ssh connections here.
    type = "ssh"

    # to connect to the instance, you must first connect to the bastion host and then to the instance itself.
    bastion_user        = "root"
    bastion_host        = "10.11.1.1"
    bastion_private_key = tls_private_key.example.private_key_pem

    user        = corellium_v1project.example.id
    host        = "proxy.enterprise.corellium.com"
    private_key = tls_private_key.example.private_key_pem
  }

  provisioner "remote-exec" {
    on_failure = continue
    inline = [
      "uptime",
    ]
  }
}

