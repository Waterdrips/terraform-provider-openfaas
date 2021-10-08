terraform {
  required_version = ">= 0.14"
  required_providers {
    openfaas = {
      version = ">= 0.0.1"
      source = "terraform.openfaas.com/openfaas/openfaas"
    }
  }
}
