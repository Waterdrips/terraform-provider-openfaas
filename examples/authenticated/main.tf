provider "openfaas" {
  uri       = "http://localhost:8080"
  user_name = "admin"
  password  = var.openfaas_provider_password
}
