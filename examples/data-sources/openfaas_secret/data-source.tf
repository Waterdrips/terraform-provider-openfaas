data "openfaas_secret" "data-secret" {
  name      = "openfaas-function"
  namespace = "openfaas-development"
}