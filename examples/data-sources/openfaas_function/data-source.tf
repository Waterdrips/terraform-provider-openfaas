data "openfaas_function" "data_function" {
  name      = "openfaas-secret"
  namespace = "openfaas-development"
}