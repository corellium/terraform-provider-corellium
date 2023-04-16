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

provider "corellium" {}

data "corellium" "example" {}
