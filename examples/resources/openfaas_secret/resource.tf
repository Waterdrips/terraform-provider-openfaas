resource "openfaas_secret" "new_secret" {
  name      = "test-function"
  namespace = "openfaas-fn"
  value     = "example secret value!"
}

