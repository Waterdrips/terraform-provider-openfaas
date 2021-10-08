provider "openfaas" {
  uri = var.openfaas_uri
  tls_insecure = true # Set to false when using a TLS endpoint
  user_name = "admin"
  password = var.openfaas_password
}