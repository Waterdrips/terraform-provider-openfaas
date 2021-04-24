terraform {
  required_version = ">= 0.14"
  required_providers {
    openfaas = {
      version = ">= 1.0"
      source = "terraform.openfaas.com/openfaas/openfaas"
      # TF-UPGRADE-TODO
      #
      # No source detected for this provider. You must add a source address
      # in the following format:
      #
      # source = "your-registry.example.com/organization/openfaas"
      #
      # For more information, see the provider source documentation:
      #
      # https://www.terraform.io/docs/configuration/providers.html#provider-source
    }
  }
}
