terraform {
  required_providers {
    hsdp = {
      source  = "philips-software/hsdp"
      version = ">= 0.42.1"
    }
    cloudfoundry = {
      source  = "cloudfoundry-community/cloudfoundry"
      version = ">= 0.50.4"
    }
  }
}
