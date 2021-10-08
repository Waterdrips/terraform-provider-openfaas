# Terraform provider for OpenFaaS

The terraform provider for [OpenFaaS](https://www.openfaas.com/)

## Documentation

Full documentation, see: // TODO


### Example Usage

Add the provider config 
```hcl
provider "openfaas" {
  uri       = "http://localhost:8080"
  user_name = "admin"
  password  = var.openfaas_provider_password
}

terraform {
  required_version = ">= 0.14"
  required_providers {
    openfaas = {
      version = ">= 0.0.1"
      source = "terraform.openfaas.com/openfaas/openfaas"
    }
  }
}

```

And then some resources
```hcl
resource "openfaas_function" "function_test" {
  name      = "test-function"
  image     = "functions/alpine:latest"
  f_process = "sha512sum"
  labels = {
    Group       = "London"
    Environment = "Test"
  }

  annotations = {
    CreatedDate = "Mon Sep  3 07:15:55 BST 2018"
  }
}
```


## Building and Installing
Since this isn't yet published to the community providers we have to manually install

### Download a release

Download and unzip the [latest release](https://github.com/Waterdrips/terraform-provider-openfaas/releases/latest).

Then, move the binary to your terraform plugins directory. [The docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
don't fully describe where this is.

* On X86_64 Mac (Not M1), it's `~/.terraform.d/plugins/terraform.openfaas.com/openfaas/openfaas/{{VERSION}}/darwin_amd64/terraform-provider-openfaas`
* On Linux, it's `~/.terraform.d/plugins/terraform.openfaas.com/openfaas/openfaas/{{VERSION}}/linux_amd64/terraform-provider-openfaas`

## Developing the Provider

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

> Note: At the moment the acceptance tests assume OpenFaaS gateway is running on http://localhost:8080 *without* 
basic authentication enabled.

```sh
$ make testacc
```

## Building the documentation

To generate docs run `tfplugindocs` (after installation)

# TODO
- [ ] Test update code
- [ ] update sdk to newer version once released
- [ ] release to hashicorp registry
- [ ] more extensive testing
- [ ] test OIDC auth
- [ ] Write tests to spin up own clusters? or makefile to make that happen?
- [ ] faas-netes/faasd provider code (work out if we support limits/reqs as faasd only support mem lim)