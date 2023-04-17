terraform {
  required_providers {
    corellium = {
      source = "aimoda/corellium"
    }
    # hashicups = {
    #   source = "hashicorp.com/edu/hashicups-pf"
    # }
  }
}

provider "corellium" {
  # placeholder token - replace with real token or use env var CORELLIUM_TOKEN
  token = "dsanjkd21en12n1"
}

data "corellium_v1ready" "example" {}

output "corellium_status_check" {
  value = data.corellium_v1ready.example
}
